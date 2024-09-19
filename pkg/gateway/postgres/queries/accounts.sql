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
select * from accounts;

-- name: UpdateAccountBalance :exec
update accounts
set balance = balance + @amount::int
where id = @id;