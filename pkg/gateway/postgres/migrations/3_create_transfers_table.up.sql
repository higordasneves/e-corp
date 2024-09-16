begin;

    create table if not exists transfers
    (
        id                     uuid        primary key ,
        account_origin_id      uuid        not null  references accounts (id),
        account_destination_id uuid        not null references accounts (id),
        amount                 bigint      not null,
        created_at             timestamptz not null,
        updated_at timestamptz not null default now()
    );

    create index on transfers (id);
    create index on transfers (account_origin_id);
    create index on transfers (account_destination_id);

    create or replace trigger tg_transfers_updated_at
        before update
        on transfers
        for each row
    execute procedure fn_trigger_updated_at();

commit;