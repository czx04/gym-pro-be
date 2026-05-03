CREATE TABLE IF NOT EXISTS social_media_assets (
    public_id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    resource_type VARCHAR(20) NOT NULL,
    secure_url TEXT,
    bytes BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'uploading',
    post_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    attached_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE SET NULL,
    CONSTRAINT check_social_media_resource_type CHECK (resource_type IN ('image', 'video')),
    CONSTRAINT check_social_media_status CHECK (status IN ('uploading', 'ready', 'attached', 'failed', 'orphaned')),
    CONSTRAINT check_social_media_bytes CHECK (bytes IS NULL OR bytes >= 0)
);

CREATE INDEX IF NOT EXISTS idx_social_media_assets_user ON social_media_assets(user_id);
CREATE INDEX IF NOT EXISTS idx_social_media_assets_status ON social_media_assets(status);
CREATE INDEX IF NOT EXISTS idx_social_media_assets_user_status ON social_media_assets(user_id, status);
CREATE INDEX IF NOT EXISTS idx_social_media_assets_expires ON social_media_assets(expires_at);

CREATE TABLE IF NOT EXISTS post_media (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    public_id TEXT NOT NULL,
    resource_type VARCHAR(20) NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (public_id) REFERENCES social_media_assets(public_id) ON DELETE CASCADE,
    CONSTRAINT check_post_media_resource_type CHECK (resource_type IN ('image', 'video')),
    CONSTRAINT check_post_media_order_index CHECK (order_index >= 0),
    CONSTRAINT unique_post_media UNIQUE (post_id, public_id)
);

CREATE INDEX IF NOT EXISTS idx_post_media_post ON post_media(post_id);
CREATE INDEX IF NOT EXISTS idx_post_media_order ON post_media(post_id, order_index);

ALTER TABLE posts
    ALTER COLUMN content_id DROP NOT NULL,
    ALTER COLUMN content_type SET DEFAULT 'general';

ALTER TABLE posts
    DROP CONSTRAINT IF EXISTS check_content_type;

ALTER TABLE posts
    ADD CONSTRAINT check_content_type CHECK (content_type IN ('general', 'workout_plan', 'meal_log'));

ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS feeling VARCHAR(100),
    ADD COLUMN IF NOT EXISTS location_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS hashtags TEXT[] NOT NULL DEFAULT '{}'::TEXT[];