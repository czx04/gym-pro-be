DROP INDEX IF EXISTS idx_user_blocks_blocked;
DROP INDEX IF EXISTS idx_user_blocks_blocker;
DROP TABLE IF EXISTS user_blocks;

DROP INDEX IF EXISTS idx_post_reports_status;
DROP INDEX IF EXISTS idx_post_reports_post;
DROP TABLE IF EXISTS post_reports;

DROP INDEX IF EXISTS idx_post_preferences_user_preference;
DROP TABLE IF EXISTS post_preferences;
