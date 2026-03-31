CREATE TABLE IF NOT EXISTS in_app_notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(32) NOT NULL,
    title TEXT NOT NULL,
    meta TEXT NOT NULL,
    post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_in_app_notifications_user_created
    ON in_app_notifications (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_in_app_notifications_user_unread
    ON in_app_notifications (user_id)
    WHERE is_read = false;
