-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS messages ADD COLUMN IF NOT EXISTS title varchar not null;
ALTER TABLE IF EXISTS messages ADD COLUMN IF NOT EXISTS author varchar not null;
ALTER TABLE IF EXISTS messages ADD COLUMN IF NOT EXISTS created_at numeric not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS messages DROP COLUMN IF EXISTS title;
ALTER TABLE IF EXISTS messages DROP COLUMN IF EXISTS author;
ALTER TABLE IF EXISTS messages DROP COLUMN IF EXISTS created_at;
-- +goose StatementEnd
