CREATE TABLE IF NOT EXISTS user_weight_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    weight_kg DECIMAL(5, 2) NOT NULL,
    measured_at TIMESTAMP WITH TIME ZONE NOT NULL,
    source VARCHAR(30) NOT NULL DEFAULT 'profile_update',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_weight_history_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_user_weight_history_weight
        CHECK (weight_kg > 0 AND weight_kg <= 500),
    CONSTRAINT check_user_weight_history_source
        CHECK (source IN ('profile_update', 'backfill_initial'))
);

CREATE INDEX IF NOT EXISTS idx_user_weight_history_user_measured_at
    ON user_weight_history(user_id, measured_at DESC);

INSERT INTO user_weight_history (user_id, weight_kg, measured_at, source)
SELECT
    u.id,
    u.weight_kg,
    COALESCE(u.updated_at, u.created_at, CURRENT_TIMESTAMP),
    'backfill_initial'
FROM users u
WHERE u.weight_kg IS NOT NULL
  AND NOT EXISTS (
      SELECT 1
      FROM user_weight_history h
      WHERE h.user_id = u.id
  );
