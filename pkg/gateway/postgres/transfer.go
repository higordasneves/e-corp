package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type transferRepo struct {
	dbPool *pgxpool.Pool
}

func NewTransferRepository(dbPool *pgxpool.Pool) repository.TransferRepo {
	return &transferRepo{dbPool: dbPool}
}

func (tRepo transferRepo) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	var db pgxtype.Querier
	db = tRepo.dbPool

	if tx := ctx.Value("dbConnection"); tx != nil {
		db = tx.(*pgxpool.Tx)
	}

	_, err := db.Exec(ctx, `INSERT INTO transfers 
		(id, account_origin_id, account_destination_id, amount, created_at)
		 VALUES ($1, $2, $3, $4, $5)`, transfer.ID.String(), transfer.AccountOriginID.String(), transfer.AccountDestinationID.String(), transfer.Amount, transfer.CreatedAt)

	if err != nil {
		return repository.NewDBError(repository.QueryRefCreateTransfer, err, repository.ErrUnexpected)
	}

	return nil
}

func (tRepo transferRepo) Transfer(ctx context.Context, ctxChan chan context.Context, errChan chan error) error {
	return PerformTransaction(ctx, ctxChan, tRepo.dbPool, errChan)
}

func (tRepo transferRepo) FetchTransfers(ctx context.Context, id vos.UUID) ([]entities.Transfer, error) {
	transferCount := tRepo.dbPool.QueryRow(ctx,
		`select count(*) as count
			from transfers
			where account_origin_id = $1`, id.String())

	var count int
	err := transferCount.Scan(&count)
	if err != nil {
		return nil, err
	}
	transferList := make([]entities.Transfer, 0, count)

	rows, err := tRepo.dbPool.Query(ctx,
		`select id
				, account_origin_id
				, account_destination_id
				, amount
				, created_at
			from transfers
			where account_origin_id = $1`, id.String())

	defer rows.Close()
	if err != nil {
		return nil, repository.NewDBError(repository.QueryRefGetTransfers, err, repository.ErrUnexpected)
	}

	for rows.Next() {
		var transfer entities.Transfer
		err = rows.Scan(&transfer.ID, &transfer.AccountOriginID, &transfer.AccountDestinationID, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return nil, repository.NewDBError(repository.QueryRefGetTransfers, err, repository.ErrUnexpected)
		}
		transferList = append(transferList, transfer)
	}
	return transferList, nil
}
