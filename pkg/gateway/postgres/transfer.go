package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func (r Repository) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	var db Querier
	db = r.dbPool

	if tx := ctx.Value("dbConnection"); tx != nil {
		db = tx.(*pgxpool.Tx)
	}

	_, err := db.Exec(ctx, `INSERT INTO transfers 
		(id, account_origin_id, account_destination_id, amount, created_at)
		 VALUES ($1, $2, $3, $4, $5)`, transfer.ID.String(), transfer.AccountOriginID.String(), transfer.AccountDestinationID.String(), transfer.Amount, transfer.CreatedAt)

	if err != nil {
		return domain.NewDBError(domain.QueryRefCreateTransfer, err, domain.ErrUnexpected)
	}

	return nil
}

func (r Repository) PerformTransaction(ctx context.Context, ctxChan chan context.Context, errChan chan error) error {
	return PerformTransaction(ctx, ctxChan, r.dbPool, errChan)
}

func (r Repository) FetchTransfers(ctx context.Context, id vos.UUID) ([]entities.Transfer, error) {
	transferCount := r.dbPool.QueryRow(ctx,
		`select count(*) as count
			from transfers
			where account_origin_id = $1`, id.String())

	var count int
	err := transferCount.Scan(&count)
	if err != nil {
		return nil, err
	}
	transferList := make([]entities.Transfer, 0, count)

	rows, err := r.dbPool.Query(ctx,
		`select id
				, account_origin_id
				, account_destination_id
				, amount
				, created_at
			from transfers
			where account_origin_id = $1`, id.String())

	defer rows.Close()
	if err != nil {
		return nil, domain.NewDBError(domain.QueryRefGetTransfers, err, domain.ErrUnexpected)
	}

	for rows.Next() {
		var transfer entities.Transfer
		err = rows.Scan(&transfer.ID, &transfer.AccountOriginID, &transfer.AccountDestinationID, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return nil, domain.NewDBError(domain.QueryRefGetTransfers, err, domain.ErrUnexpected)
		}
		transferList = append(transferList, transfer)
	}
	return transferList, nil
}
