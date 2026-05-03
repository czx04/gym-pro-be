UPDATE in_app_notifications
SET meta = meta || '|' || related_post_id::text
WHERE related_post_id IS NOT NULL
  AND strpos(meta, '|') = 0;
