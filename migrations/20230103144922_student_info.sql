-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS student_infos (
	id uuid NOT NULL,
	name text NULL,
	roll_number text NULL,
    age numeric NULL, 
	course_id uuid NULL,
	marks_id uuid NULL,
	CONSTRAINT student_infos_pkey PRIMARY KEY (id),
    CONSTRAINT student_infos_course_id_fkey FOREIGN KEY (course_id) REFERENCES public.course_infos(id),
    CONSTRAINT student_infos_marks_id_fkey FOREIGN KEY(marks_id) REFERENCES public.student_marks(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS student_infos DROP CONSTRAINT IF EXISTS student_infos_course_id_fkey;
ALTER TABLE IF EXISTS student_infos DROP CONSTRAINT IF EXISTS student_infos_marks_id_fkey;
DROP TABLE IF EXISTS student_infos;
-- +goose StatementEnd
