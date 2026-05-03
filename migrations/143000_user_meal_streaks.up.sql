CREATE TABLE IF NOT EXISTS user_meal_streaks (
    user_id UUID PRIMARY KEY,
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_meal_streaks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_current_streak CHECK (current_streak >= 0),
    CONSTRAINT check_longest_streak CHECK (longest_streak >= 0)
);

CREATE INDEX IF NOT EXISTS idx_user_meal_streaks_updated ON user_meal_streaks(updated_at);
