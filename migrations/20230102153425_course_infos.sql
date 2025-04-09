-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS course_infos (
	id uuid NOT NULL,
	course_name text NULL,
	CONSTRAINT course_infos_pkey PRIMARY KEY (id)
);

INSERT INTO public.course_infos
(id, course_name)
VALUES('285aa383-53b5-4e26-bf3a-c76b3245c617', 'No Course');

INSERT INTO course_infos (id, course_name) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Introduction to Computer Science'),
('123e4567-e89b-12d3-a456-426614174001', 'Data Structures and Algorithms'),
('987f6543-b21c-34d2-c789-123456789abc', 'Database Management Systems'),
('456e7890-c12d-45e6-b789-987654321def', 'Operating Systems'),
('c88fef09-97b3-4c97-b05d-b6bd72288d51', 'Artificial Intelligence'),
('56c8c014-e5fa-40f5-9ac1-42bda3c884cb', 'Machine Learning'),
('8ae462c5-d418-41b7-b675-63affb4fbe16', 'Software Engineering'),
('745132ae-5bcf-4c5f-8ca2-d926b16ed1bd', 'Computer Networks'),
('42803923-1d87-4abd-9285-c2a50a677d22', 'Web Development'),
('5873fb7a-26d7-47e5-b5d1-e369255fc878', 'Cloud Computing');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_infos;
-- +goose StatementEnd
