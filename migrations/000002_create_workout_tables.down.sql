DROP INDEX IF EXISTS idx_workout_session_exercises_exercise;
DROP INDEX IF EXISTS idx_workout_session_exercises_session;
DROP TABLE IF EXISTS workout_session_exercises;

DROP INDEX IF EXISTS idx_workout_sessions_user_started;
DROP INDEX IF EXISTS idx_workout_sessions_started;
DROP INDEX IF EXISTS idx_workout_sessions_user;
DROP TABLE IF EXISTS workout_sessions;

DROP INDEX IF EXISTS idx_workout_schedules_completed;
DROP INDEX IF EXISTS idx_workout_schedules_user_date;
DROP INDEX IF EXISTS idx_workout_schedules_date;
DROP INDEX IF EXISTS idx_workout_schedules_user;
DROP TABLE IF EXISTS workout_schedules;

DROP INDEX IF EXISTS idx_workout_plan_exercises_exercise;
DROP INDEX IF EXISTS idx_workout_plan_exercises_plan;
DROP TABLE IF EXISTS workout_plan_exercises;

DROP INDEX IF EXISTS idx_workout_plans_public;
DROP INDEX IF EXISTS idx_workout_plans_user;
DROP TABLE IF EXISTS workout_plans;

DROP INDEX IF EXISTS idx_exercises_name;
DROP INDEX IF EXISTS idx_exercises_is_active;
DROP INDEX IF EXISTS idx_exercises_difficulty;
DROP INDEX IF EXISTS idx_exercises_category;
DROP TABLE IF EXISTS exercises;
