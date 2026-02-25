DROP INDEX IF EXISTS idx_meal_log_items_recipe;
DROP INDEX IF EXISTS idx_meal_log_items_food;
DROP INDEX IF EXISTS idx_meal_log_items_log;
DROP TABLE IF EXISTS meal_log_items;

DROP INDEX IF EXISTS idx_meal_logs_user_date_time;
DROP INDEX IF EXISTS idx_meal_logs_user_date;
DROP INDEX IF EXISTS idx_meal_logs_date;
DROP INDEX IF EXISTS idx_meal_logs_user;
DROP TABLE IF EXISTS meal_logs;

DROP INDEX IF EXISTS idx_recipe_foods_food;
DROP INDEX IF EXISTS idx_recipe_foods_recipe;
DROP TABLE IF EXISTS recipe_foods;

DROP INDEX IF EXISTS idx_recipes_name;
DROP INDEX IF EXISTS idx_recipes_public;
DROP INDEX IF EXISTS idx_recipes_user;
DROP TABLE IF EXISTS recipes;

DROP INDEX IF EXISTS idx_foods_barcode;
DROP INDEX IF EXISTS idx_foods_user;
DROP INDEX IF EXISTS idx_foods_is_system;
DROP INDEX IF EXISTS idx_foods_category;
DROP INDEX IF EXISTS idx_foods_name;
DROP TABLE IF EXISTS foods;
