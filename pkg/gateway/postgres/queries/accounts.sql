-- name: InsertAccount :exec
insert into accounts (id, document_number, name, secret, balance, created_at)
values (@id, @document_number, @name, @secret, @balance, @created_at);

-- name: GetAccount :one
select *
from accounts
where id = @id;

-- name: GetAccountByDocument :one
select *
from accounts
where document_number = @document_number;

-- name: ListAccounts :many
select * from accounts a
where id = any(@ids::uuid[])
    and (@last_fetched_id::uuid = '00000000-0000-0000-0000-000000000000'  or id < @last_fetched_id::uuid)
order by a.id desc
limit @page_size;

-- name: UpdateAccountBalance :exec
update accounts
set balance = balance + @amount::int
where id = @id;