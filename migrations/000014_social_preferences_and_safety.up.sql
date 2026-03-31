CREATE TABLE IF NOT EXISTS post_preferences (
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    preference VARCHAR(30) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT check_post_preference_type CHECK (preference IN ('interested', 'not_interested'))
);

CREATE INDEX IF NOT EXISTS idx_post_preferences_user_preference
    ON post_preferences(user_id, preference);

CREATE TABLE IF NOT EXISTS post_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    reporter_id UUID NOT NULL,
    reason VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_post_report UNIQUE (post_id, reporter_id),
    CONSTRAINT check_post_report_reason CHECK (reason IN ('spam', 'harassment', 'misinformation', 'nudity', 'violence', 'other')),
    CONSTRAINT check_post_report_status CHECK (status IN ('pending', 'reviewing', 'resolved', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_post_reports_post ON post_reports(post_id);
CREATE INDEX IF NOT EXISTS idx_post_reports_status ON post_reports(status);

CREATE TABLE IF NOT EXISTS user_blocks (
    blocker_id UUID NOT NULL,
    blocked_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (blocker_id, blocked_id),
    FOREIGN KEY (blocker_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blocked_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_no_self_block CHECK (blocker_id <> blocked_id)
);

CREATE INDEX IF NOT EXISTS idx_user_blocks_blocker ON user_blocks(blocker_id);
CREATE INDEX IF NOT EXISTS idx_user_blocks_blocked ON user_blocks(blocked_id);
