-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS accounts
ADD COLUMN IF NOT EXISTS verified BOOLEAN DEFAULT FALSE;

UPDATE accounts
SET verified = TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS accounts
DROP COLUMN IF EXISTS verified
-- +goose StatementEnd
