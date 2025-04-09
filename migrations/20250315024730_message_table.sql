-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages(
    id uuid not null,
    account_id uuid not null,
    messages varchar null,
    is_read boolean,
    CONSTRAINT messages_pkey PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS messages DROP CONSTRAINT IF EXISTS messages_pkey;
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
