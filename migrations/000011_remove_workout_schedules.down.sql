-- Restore workout_schedules and optional link from workout_sessions
CREATE TABLE IF NOT EXISTS workout_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_plan_id UUID NOT NULL,
    user_id UUID NOT NULL,
    scheduled_date DATE NOT NULL,
    scheduled_time TIME,
    recurrence_rule VARCHAR(100),
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (workout_plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_workout_schedules_user ON workout_schedules(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_schedules_date ON workout_schedules(scheduled_date);
CREATE INDEX IF NOT EXISTS idx_workout_schedules_user_date ON workout_schedules(user_id, scheduled_date);
CREATE INDEX IF NOT EXISTS idx_workout_schedules_completed ON workout_schedules(is_completed);

ALTER TABLE workout_sessions
    ADD COLUMN IF NOT EXISTS workout_schedule_id UUID;

ALTER TABLE workout_sessions
    DROP CONSTRAINT IF EXISTS workout_sessions_workout_schedule_id_fkey;

ALTER TABLE workout_sessions
    ADD CONSTRAINT workout_sessions_workout_schedule_id_fkey
    FOREIGN KEY (workout_schedule_id)
    REFERENCES workout_schedules(id)
    ON DELETE SET NULL;
