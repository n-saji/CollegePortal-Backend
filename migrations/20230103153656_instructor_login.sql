-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS instructor_logins (
    id uuid not null,
    email_id text null,
    password text null,
    CONSTRAINT instructor_logins_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE If EXISTS instructor_logins;
-- +goose StatementEnd
