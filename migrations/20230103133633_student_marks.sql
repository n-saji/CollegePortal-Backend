-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS student_marks (
    id uuid NOT NULL,
    student_id uuid NULL,
    course_id uuid NULL,
    course_name text NULL,
    marks numeric NULL,
    grade text NULL,
    CONSTRAINT student_marks_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_marks;
-- +goose StatementEnd
