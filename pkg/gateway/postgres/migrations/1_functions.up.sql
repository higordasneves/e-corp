begin;

    create or replace function fn_trigger_updated_at()
        returns trigger
        language plpgsql
    as
    $$
    begin
        new.updated_at = now();
    return new;
    end;
    $$;

commit;