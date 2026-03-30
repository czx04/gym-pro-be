DROP INDEX IF EXISTS idx_foods_embedding;
ALTER TABLE foods DROP COLUMN IF EXISTS embedding;
