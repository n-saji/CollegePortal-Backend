-- +goose Up
-- +goose StatementBegin
Alter table if exists instructor_details 
add column if not exists info jsonb default '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Alter table if exists instructor_details 
drop column if exists info;
-- +goose StatementEnd
