-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_default_status_grade()
RETURNS TRIGGER AS $my_table$
BEGIN
    UPDATE student_marks
    SET grade = 'nil'
    WHERE grade = '';
RETURN NEW;
END;
$my_table$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER setting_default_grade
after
INSERT ON student_infos FOR EACH ROW EXECUTE FUNCTION set_default_status_grade();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS setting_default_grade on student_infos;
-- +goose StatementEnd
