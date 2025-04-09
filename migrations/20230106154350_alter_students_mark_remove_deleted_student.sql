-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION delete_student_marks()
RETURNS TRIGGER AS $my_table$
BEGIN
   DELETE FROM public.student_marks
    WHERE student_id::text not in (SELECT id::text from  public.student_infos);
RETURN NEW;
END;
$my_table$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER deleting_student_mark_trigger
after
delete ON student_infos FOR EACH ROW EXECUTE PROCEDURE delete_student_marks();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS deleting_student_mark_trigger on student_infos;
-- +goose StatementEnd
