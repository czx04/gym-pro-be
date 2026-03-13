CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    name    TEXT NOT NULL,
    applied_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE schema_migrations IS 'Track applied migrations for auto-migrate service';
