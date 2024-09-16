begin;

    create table accounts
    (
        id              uuid        PRIMARY KEY,
        document_number text        not null,
        name            text        not null,
        secret          text        not null,
        balance         bigint      not null default 0 check (balance >= 0 ),
        created_at      timestamptz not null,
        updated_at timestamptz not null default now()
    );

    create unique index on accounts (document_number);

    create or replace trigger tg_accounts_updated_at
    before update
    on accounts
    for each row
    execute procedure fn_trigger_updated_at();

commit;