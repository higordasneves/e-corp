// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: accounts.sql

package sqlc

import (
	"context"
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

const GetAccount = `-- name: GetAccount :one
select id, document_number, name, secret, balance, created_at, updated_at
from accounts
where id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id uuid.UUID) (Account, error) {
	row := q.db.QueryRow(ctx, GetAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.DocumentNumber,
		&i.Name,
		&i.Secret,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetAccountByDocument = `-- name: GetAccountByDocument :one
select id, document_number, name, secret, balance, created_at, updated_at
from accounts
where document_number = $1
`

func (q *Queries) GetAccountByDocument(ctx context.Context, documentNumber string) (Account, error) {
	row := q.db.QueryRow(ctx, GetAccountByDocument, documentNumber)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.DocumentNumber,
		&i.Name,
		&i.Secret,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const InsertAccount = `-- name: InsertAccount :exec
insert into accounts (id, document_number, name, secret, balance, created_at)
values ($1, $2, $3, $4, $5, $6)
`

type InsertAccountParams struct {
	ID             uuid.UUID
	DocumentNumber string
	Name           string
	Secret         string
	Balance        int64
	CreatedAt      time.Time
}

func (q *Queries) InsertAccount(ctx context.Context, arg InsertAccountParams) error {
	_, err := q.db.Exec(ctx, InsertAccount,
		arg.ID,
		arg.DocumentNumber,
		arg.Name,
		arg.Secret,
		arg.Balance,
		arg.CreatedAt,
	)
	return err
}

const ListAccounts = `-- name: ListAccounts :many
select id, document_number, name, secret, balance, created_at, updated_at from accounts a
where id = any($1::uuid[])
    and ($2::uuid = '00000000-0000-0000-0000-000000000000'  or id < $2::uuid)
order by a.id desc
limit $3
`

type ListAccountsParams struct {
	Ids           []uuid.UUID
	LastFetchedID uuid.UUID
	PageSize      int32
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, ListAccounts, arg.Ids, arg.LastFetchedID, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.DocumentNumber,
			&i.Name,
			&i.Secret,
			&i.Balance,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateAccountBalance = `-- name: UpdateAccountBalance :exec
update accounts
set balance = balance + $1::int
where id = $2
`

type UpdateAccountBalanceParams struct {
	Amount int32
	ID     uuid.UUID
}

func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) error {
	_, err := q.db.Exec(ctx, UpdateAccountBalance, arg.Amount, arg.ID)
	return err
}
