ALTER TABLE workout_session_sets
    ADD COLUMN IF NOT EXISTS rest_secs INTEGER;

ALTER TABLE workout_session_sets
    DROP CONSTRAINT IF EXISTS check_workout_session_sets_rest_secs;

ALTER TABLE workout_session_sets
    ADD CONSTRAINT check_workout_session_sets_rest_secs
    CHECK (rest_secs IS NULL OR (rest_secs >= 0 AND rest_secs <= 1800));

COMMENT ON COLUMN workout_session_sets.rest_secs IS 'Actual rest time between sets in seconds, tracked during session updates';
