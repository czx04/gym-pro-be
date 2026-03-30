ALTER TABLE in_app_notifications
    ADD COLUMN IF NOT EXISTS related_post_id UUID;

UPDATE in_app_notifications
SET related_post_id = post_id
WHERE related_post_id IS NULL
  AND post_id IS NOT NULL;
