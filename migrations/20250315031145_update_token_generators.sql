-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS token_generators ADD COLUMN IF NOT EXISTS account_id uuid null;
ALTER TABLE IF EXISTS token_generators ADD CONSTRAINT token_generators_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS token_generators DROP CONSTRAINT IF EXISTS token_generators_fkey;
ALTER TABLE IF EXISTS token_generators DROP COLUMN IF EXISTS account_id;
-- +goose StatementEnd
