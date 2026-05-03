DELETE FROM posts WHERE content_type = 'general';

ALTER TABLE posts
	DROP COLUMN IF EXISTS hashtags,
	DROP COLUMN IF EXISTS location_name,
	DROP COLUMN IF EXISTS feeling;

ALTER TABLE posts
	ALTER COLUMN content_type DROP DEFAULT;

ALTER TABLE posts
	DROP CONSTRAINT IF EXISTS check_content_type;

ALTER TABLE posts
	ADD CONSTRAINT check_content_type CHECK (content_type IN ('workout_plan', 'meal_log'));

ALTER TABLE posts
	ALTER COLUMN content_id SET NOT NULL;

DROP INDEX IF EXISTS idx_post_media_order;
DROP INDEX IF EXISTS idx_post_media_post;
DROP TABLE IF EXISTS post_media;

DROP INDEX IF EXISTS idx_social_media_assets_expires;
DROP INDEX IF EXISTS idx_social_media_assets_user_status;
DROP INDEX IF EXISTS idx_social_media_assets_status;
DROP INDEX IF EXISTS idx_social_media_assets_user;
DROP TABLE IF EXISTS social_media_assets;