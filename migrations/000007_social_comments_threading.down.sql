DROP INDEX IF EXISTS idx_comments_parent_created;
DROP INDEX IF EXISTS idx_comments_post_parent_created;

ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS comments_parent_comment_fk;

ALTER TABLE comments
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS reply_count,
    DROP COLUMN IF EXISTS parent_comment_id;