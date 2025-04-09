-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS token_generators (
    token uuid not null,
    valid_from numeric null,
    valid_till numeric null,
    is_valid boolean null,
    CONSTRAINT token_generators_pkey PRIMARY KEY (token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE If EXISTS token_generators;
-- +goose StatementEnd
