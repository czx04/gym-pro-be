-- Add calendar/tracking fields to workout_sessions: scheduled_date, status, started_at nullable
ALTER TABLE workout_sessions
    ADD COLUMN IF NOT EXISTS scheduled_date DATE,
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'scheduled',
    ALTER COLUMN started_at DROP NOT NULL;

COMMENT ON COLUMN workout_sessions.scheduled_date IS 'Ngày đã lên lịch (cho calendar)';
COMMENT ON COLUMN workout_sessions.status IS 'scheduled | in_progress | completed';

UPDATE workout_sessions SET status = 'completed' WHERE completed_at IS NOT NULL;
UPDATE workout_sessions SET status = 'in_progress' WHERE completed_at IS NULL AND started_at IS NOT NULL;
UPDATE workout_sessions SET status = 'scheduled' WHERE status IS NULL OR status = '';

ALTER TABLE workout_sessions
    ADD CONSTRAINT check_session_status CHECK (status IN ('scheduled', 'in_progress', 'completed'));

CREATE INDEX IF NOT EXISTS idx_workout_sessions_scheduled_date ON workout_sessions(scheduled_date);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_user_scheduled ON workout_sessions(user_id, scheduled_date) WHERE scheduled_date IS NOT NULL;

-- Bảng lưu từng set khi tracking: reps, khối lượng tạ (nullable)
CREATE TABLE IF NOT EXISTS workout_session_sets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workout_session_exercise_id UUID NOT NULL,
    set_index INTEGER NOT NULL,
    reps INTEGER,
    weight_kg DECIMAL(6, 2) NULL,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (workout_session_exercise_id) REFERENCES workout_session_exercises(id) ON DELETE CASCADE,
    CONSTRAINT unique_session_exercise_set UNIQUE (workout_session_exercise_id, set_index)
);

CREATE INDEX IF NOT EXISTS idx_workout_session_sets_exercise ON workout_session_sets(workout_session_exercise_id);

COMMENT ON COLUMN workout_session_sets.weight_kg IS 'Khối lượng tạ (kg), nullable cho bài không dùng tạ';
