DROP INDEX IF EXISTS idx_workout_session_sets_exercise;
DROP TABLE IF EXISTS workout_session_sets;

DROP INDEX IF EXISTS idx_workout_sessions_user_scheduled;
DROP INDEX IF EXISTS idx_workout_sessions_scheduled_date;
ALTER TABLE workout_sessions DROP CONSTRAINT IF EXISTS check_session_status;
ALTER TABLE workout_sessions DROP COLUMN IF EXISTS status;
ALTER TABLE workout_sessions DROP COLUMN IF EXISTS scheduled_date;
ALTER TABLE workout_sessions ALTER COLUMN started_at SET NOT NULL;
