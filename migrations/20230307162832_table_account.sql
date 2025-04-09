-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts(
    id uuid not null,
    name varchar null,
    info jsonb,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS accounts_pkey;
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
