-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS instructor_details (
    id uuid not null,
    instructor_code text null,
    instructor_name text null,
    department text null,
    course_id uuid null,
    CONSTRAINT instructor_details_pkey PRIMARY KEY (id),
	CONSTRAINT instructor_details_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.course_infos(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS instructor_details DROP CONSTRAINT IF EXISTS instructor_details_course_id_fkey;
DROP TABLE IF EXISTS instructor_details;
-- +goose StatementEnd
