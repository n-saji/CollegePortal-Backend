-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS otps (
    id uuid PRIMARY KEY NOT NULL ,
    account_id uuid NOT NULL,
    email_id VARCHAR(255) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    created_at numeric DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP),
    expires_at numeric NOT NULL,
    is_used BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS otps;
-- +goose StatementEnd
