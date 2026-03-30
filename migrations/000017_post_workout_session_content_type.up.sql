ALTER TABLE posts DROP CONSTRAINT IF EXISTS check_content_type;
ALTER TABLE posts ADD CONSTRAINT check_content_type CHECK (content_type IN ('general', 'workout_plan', 'meal_log', 'workout_session'));
