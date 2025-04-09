-- +goose Up
-- +goose StatementBegin
UPDATE accounts set type = 'instructor';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
