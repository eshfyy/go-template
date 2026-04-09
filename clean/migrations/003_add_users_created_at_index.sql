-- +goose Up
CREATE INDEX CONCURRENTLY idx_users_created_at ON users (created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_users_created_at;
