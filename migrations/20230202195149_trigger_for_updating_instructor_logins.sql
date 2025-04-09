-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION delete_instructor_login()
RETURNS TRIGGER AS $my_table$
BEGIN
   DELETE FROM public.instructor_logins
    WHERE id::text not in (SELECT id::text from  public.instructor_details);
RETURN NEW;
END;
$my_table$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER delete_instructor_login_trigger
after
delete ON instructor_details FOR EACH ROW EXECUTE PROCEDURE delete_instructor_login();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS delete_instructor_login_trigger on instructor_details;
-- +goose StatementEnd
