ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_posts_deleted_at
    ON posts(deleted_at);
