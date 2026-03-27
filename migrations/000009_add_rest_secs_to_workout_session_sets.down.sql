ALTER TABLE workout_session_sets
    DROP CONSTRAINT IF EXISTS check_workout_session_sets_rest_secs;

ALTER TABLE workout_session_sets
    DROP COLUMN IF EXISTS rest_secs;
