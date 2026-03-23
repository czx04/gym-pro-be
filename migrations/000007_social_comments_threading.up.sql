ALTER TABLE comments
    ADD COLUMN IF NOT EXISTS parent_comment_id UUID,
    ADD COLUMN IF NOT EXISTS reply_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS comments_parent_comment_fk;

ALTER TABLE comments
    ADD CONSTRAINT comments_parent_comment_fk
        FOREIGN KEY (parent_comment_id) REFERENCES comments(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_comments_post_parent_created
    ON comments(post_id, parent_comment_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_comments_parent_created
    ON comments(parent_comment_id, created_at ASC)
    WHERE deleted_at IS NULL;