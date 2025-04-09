-- +goose Up
-- +goose StatementBegin
-- CREATE OR REPLACE FUNCTION delete_account_ids_id()
-- RETURNS TRIGGER AS $my_table$
-- BEGIN
--    DELETE FROM public.accounts
--     WHERE id not in (SELECT id from  public.instructor_details);
-- RETURN NEW;
-- END;
-- $my_table$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_account_ids()
RETURNS TRIGGER AS $my_table$
BEGIN
   DELETE FROM public.accounts
    WHERE id::text not in 
        (SELECT si.id::text from  public.student_infos si 
        UNION 
        SELECT ids.id::text from  public.instructor_details ids);
RETURN NEW;
END;
$my_table$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER delete_account_ids_id_trigger
after
delete ON public.instructor_details FOR EACH ROW EXECUTE FUNCTION delete_account_ids();

CREATE OR REPLACE TRIGGER delete_account_ids_si_trigger
after
delete ON  public.student_infos FOR EACH ROW EXECUTE FUNCTION delete_account_ids();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS delete_account_ids_id_trigger on public.accounts;
DROP TRIGGER IF EXISTS delete_account_ids_si_trigger on public.accounts;
-- +goose StatementEnd