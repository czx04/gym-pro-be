CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE foods ADD COLUMN IF NOT EXISTS embedding vector(768);

CREATE INDEX IF NOT EXISTS idx_foods_embedding ON foods USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
