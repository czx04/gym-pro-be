-- Seed data for testing weekly AI summary
-- Range: 2026-03-21 -> 2026-03-28
-- Target user: c32698e0-b200-4878-b411-9d78e879cc4f
--
-- Usage (example):
-- psql "$DATABASE_URL" -f scripts/seed_weekly_ai_summary_2026_03_21_2026_03_28.sql

BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM users
        WHERE id = 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid
    ) THEN
        RAISE EXCEPTION 'User % does not exist', 'c32698e0-b200-4878-b411-9d78e879cc4f';
    END IF;
END $$;

-- 1) Ensure deterministic exercises for this seed.
INSERT INTO exercises (
    id, name, description, category, muscle_groups, equipment_needed,
    difficulty_level, calories_per_minute, is_active, created_at, updated_at
) VALUES
    (
        'f0000000-0000-0000-0000-000000000001'::uuid,
        'Barbell Bench Press',
        'Seed exercise for AI summary test',
        'strength',
        '["chest", "triceps", "shoulders"]'::jsonb,
        '["barbell", "bench"]'::jsonb,
        'intermediate',
        8.5,
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'f0000000-0000-0000-0000-000000000002'::uuid,
        'Back Squat',
        'Seed exercise for AI summary test',
        'strength',
        '["quads", "glutes", "core"]'::jsonb,
        '["barbell", "rack"]'::jsonb,
        'intermediate',
        9.2,
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'f0000000-0000-0000-0000-000000000003'::uuid,
        'Conventional Deadlift',
        'Seed exercise for AI summary test',
        'strength',
        '["hamstrings", "back", "glutes"]'::jsonb,
        '["barbell"]'::jsonb,
        'advanced',
        10.1,
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'f0000000-0000-0000-0000-000000000004'::uuid,
        'Seated Dumbbell Shoulder Press',
        'Seed exercise for AI summary test',
        'strength',
        '["shoulders", "triceps"]'::jsonb,
        '["dumbbell", "bench"]'::jsonb,
        'beginner',
        7.4,
        TRUE,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    muscle_groups = EXCLUDED.muscle_groups,
    equipment_needed = EXCLUDED.equipment_needed,
    difficulty_level = EXCLUDED.difficulty_level,
    calories_per_minute = EXCLUDED.calories_per_minute,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- 2) Create/update a fixed workout plan for this user.
INSERT INTO workout_plans (
    id, user_id, title, description, difficulty_level,
    estimated_duration_mins, estimated_calories, is_template, is_public, created_at, updated_at
) VALUES (
    'f1000000-0000-0000-0000-000000000001'::uuid,
    'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid,
    'AI Summary Seed Plan',
    'Seed plan for weekly AI summary testing',
    'intermediate',
    60,
    450,
    FALSE,
    FALSE,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    difficulty_level = EXCLUDED.difficulty_level,
    estimated_duration_mins = EXCLUDED.estimated_duration_mins,
    estimated_calories = EXCLUDED.estimated_calories,
    updated_at = NOW();

DELETE FROM workout_plan_exercises
WHERE workout_plan_id = 'f1000000-0000-0000-0000-000000000001'::uuid;

INSERT INTO workout_plan_exercises (
    id, workout_plan_id, exercise_id, "order", sets, reps, duration_secs, rest_secs, notes
) VALUES
    (
        uuid_generate_v4(),
        'f1000000-0000-0000-0000-000000000001'::uuid,
        'f0000000-0000-0000-0000-000000000001'::uuid,
        1,
        4,
        10,
        120,
        90,
        'Main push lift'
    ),
    (
        uuid_generate_v4(),
        'f1000000-0000-0000-0000-000000000001'::uuid,
        'f0000000-0000-0000-0000-000000000002'::uuid,
        2,
        4,
        8,
        150,
        120,
        'Main squat lift'
    ),
    (
        uuid_generate_v4(),
        'f1000000-0000-0000-0000-000000000001'::uuid,
        'f0000000-0000-0000-0000-000000000003'::uuid,
        3,
        3,
        6,
        180,
        150,
        'Main pull lift'
    ),
    (
        uuid_generate_v4(),
        'f1000000-0000-0000-0000-000000000001'::uuid,
        'f0000000-0000-0000-0000-000000000004'::uuid,
        4,
        3,
        12,
        90,
        75,
        'Accessory press'
    );

-- 3) Cleanup existing sessions in requested range (idempotent reseed).
DELETE FROM workout_sessions
WHERE user_id = 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid
  AND scheduled_date BETWEEN DATE '2026-03-21' AND DATE '2026-03-28';

-- 4) Seed sessions for 8 days.
WITH seed_days AS (
    SELECT *
    FROM (VALUES
        (DATE '2026-03-21', TIME '06:30', 58, 360, 'energetic', 3, 'Strong start'),
        (DATE '2026-03-22', TIME '06:35', 62, 390, 'happy', 3, 'Good form'),
        (DATE '2026-03-23', TIME '06:40', 55, 340, 'neutral', 4, 'Slightly heavy day'),
        (DATE '2026-03-24', TIME '06:45', 64, 410, 'happy', 4, 'Progressive overload'),
        (DATE '2026-03-25', TIME '06:50', 57, 355, 'neutral', 3, 'Recovery-focused'),
        (DATE '2026-03-26', TIME '06:55', 66, 425, 'energetic', 4, 'Top performance'),
        (DATE '2026-03-27', TIME '07:00', 60, 370, 'happy', 3, 'Solid consistency'),
        (DATE '2026-03-28', TIME '07:05', 52, 320, 'tired', 4, 'Fatigue check day')
    ) AS v(scheduled_date, start_time, duration_mins, calories, mood, difficulty_rating, notes)
),
inserted_sessions AS (
    INSERT INTO workout_sessions (
        id,
        user_id,
        workout_plan_id,
        scheduled_date,
        status,
        started_at,
        completed_at,
        duration_mins,
        total_calories_burned,
        notes,
        mood,
        difficulty_rating,
        created_at,
        updated_at
    )
    SELECT
        uuid_generate_v4(),
        'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid,
        'f1000000-0000-0000-0000-000000000001'::uuid,
        d.scheduled_date,
        'completed',
        (d.scheduled_date::timestamp + d.start_time),
        (d.scheduled_date::timestamp + d.start_time + make_interval(mins => d.duration_mins)),
        d.duration_mins,
        d.calories,
        d.notes,
        d.mood,
        d.difficulty_rating,
        NOW(),
        NOW()
    FROM seed_days d
    RETURNING id
)
INSERT INTO workout_session_exercises (
    id,
    workout_session_id,
    exercise_id,
    "order",
    target_sets,
    target_reps,
    duration_secs,
    notes,
    skipped
)
SELECT
    uuid_generate_v4(),
    s.id,
    p.exercise_id,
    p."order",
    p.sets,
    p.reps,
    p.duration_secs,
    'Seed tracked exercise',
    FALSE
FROM inserted_sessions s
CROSS JOIN LATERAL (
    SELECT exercise_id, "order", sets, reps, duration_secs
    FROM workout_plan_exercises
    WHERE workout_plan_id = 'f1000000-0000-0000-0000-000000000001'::uuid
) p;

-- 5) Seed set-level data (reps, weight, rest_secs, completed_at) for each session exercise.
INSERT INTO workout_session_sets (
    id,
    workout_session_exercise_id,
    set_index,
    reps,
    weight_kg,
    rest_secs,
    completed,
    completed_at,
    created_at,
    updated_at
)
SELECT
    uuid_generate_v4(),
    wse.id,
    gs.set_idx,
    GREATEST(3, COALESCE(wse.target_reps, 10) - (gs.set_idx - 1)),
    ROUND((
        CASE wse."order"
            WHEN 1 THEN 42.0
            WHEN 2 THEN 60.0
            WHEN 3 THEN 85.0
            WHEN 4 THEN 22.0
            ELSE 30.0
        END
        + ((EXTRACT(DAY FROM ws.scheduled_date)::int - 21) * 0.7)
        + ((gs.set_idx - 1) * 1.5)
    )::numeric, 2),
    LEAST(1800, 65 + (gs.set_idx * 8) + (wse."order" * 5)),
    TRUE,
    COALESCE(ws.started_at, ws.created_at) + make_interval(mins => (wse."order" * 7) + gs.set_idx),
    NOW(),
    NOW()
FROM workout_session_exercises wse
JOIN workout_sessions ws ON ws.id = wse.workout_session_id
CROSS JOIN LATERAL generate_series(1, COALESCE(wse.target_sets, 3)) AS gs(set_idx)
WHERE ws.user_id = 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid
  AND ws.scheduled_date BETWEEN DATE '2026-03-21' AND DATE '2026-03-28';

-- 6) Seed weight history in the same date range for body_weight_trend.
DELETE FROM user_weight_history
WHERE user_id = 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid
  AND measured_at::date BETWEEN DATE '2026-03-21' AND DATE '2026-03-28'
  AND source IN ('profile_update', 'backfill_initial');

INSERT INTO user_weight_history (id, user_id, weight_kg, measured_at, source, created_at)
VALUES
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 72.40, TIMESTAMPTZ '2026-03-21 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 72.10, TIMESTAMPTZ '2026-03-22 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 71.95, TIMESTAMPTZ '2026-03-23 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 71.80, TIMESTAMPTZ '2026-03-24 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 71.70, TIMESTAMPTZ '2026-03-25 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 71.60, TIMESTAMPTZ '2026-03-26 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 71.55, TIMESTAMPTZ '2026-03-27 06:00:00+07', 'profile_update', NOW()),
    (uuid_generate_v4(), 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid, 71.45, TIMESTAMPTZ '2026-03-28 06:00:00+07', 'profile_update', NOW());

COMMIT;

-- Quick sanity checks:
-- SELECT scheduled_date, mood, difficulty_rating, duration_mins, total_calories_burned
-- FROM workout_sessions
-- WHERE user_id = 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid
--   AND scheduled_date BETWEEN DATE '2026-03-21' AND DATE '2026-03-28'
-- ORDER BY scheduled_date;
--
-- SELECT measured_at::date, weight_kg
-- FROM user_weight_history
-- WHERE user_id = 'c32698e0-b200-4878-b411-9d78e879cc4f'::uuid
--   AND measured_at::date BETWEEN DATE '2026-03-21' AND DATE '2026-03-28'
-- ORDER BY measured_at;
