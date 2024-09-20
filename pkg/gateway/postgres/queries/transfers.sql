-- name: InsertTransfer :exec
insert into transfers(id, account_origin_id, account_destination_id, amount, created_at)
values (@id, @account_origin_id, @account_destination_id, @amount, @created_at);

-- name: ListAccountSentTransfers :many
select *
from transfers
where account_origin_id = @account_origin_id;