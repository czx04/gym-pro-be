CREATE TABLE IF NOT EXISTS user_daily_steps (
    user_id UUID NOT NULL,
    date DATE NOT NULL,
    steps INTEGER NOT NULL DEFAULT 0,
    source VARCHAR(30) NOT NULL DEFAULT 'apple_health',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_daily_steps_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_user_daily_steps_steps CHECK (steps >= 0),
    CONSTRAINT check_user_daily_steps_source CHECK (source IN ('apple_health'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_daily_steps_user_date_source
    ON user_daily_steps(user_id, date, source);

CREATE INDEX IF NOT EXISTS idx_user_daily_steps_user_date
    ON user_daily_steps(user_id, date DESC);

