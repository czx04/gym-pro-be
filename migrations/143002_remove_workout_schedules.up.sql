-- Consolidate workout calendar/tracking to workout_sessions as single source of truth
ALTER TABLE workout_sessions
    DROP CONSTRAINT IF EXISTS workout_sessions_workout_schedule_id_fkey;

ALTER TABLE workout_sessions
    DROP COLUMN IF EXISTS workout_schedule_id;

DROP TABLE IF EXISTS workout_schedules;
