CREATE TABLE IF NOT EXISTS comment_media (
    comment_id UUID NOT NULL,
    public_id VARCHAR(255) NOT NULL,
    resource_type VARCHAR(20) NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (comment_id, public_id),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (public_id) REFERENCES social_media_assets(public_id) ON DELETE CASCADE,
    CONSTRAINT check_comment_media_resource_type CHECK (resource_type IN ('image')),
    CONSTRAINT check_comment_media_order_index CHECK (order_index >= 0)
);

CREATE INDEX IF NOT EXISTS idx_comment_media_comment ON comment_media(comment_id);
CREATE INDEX IF NOT EXISTS idx_comment_media_order ON comment_media(comment_id, order_index);
