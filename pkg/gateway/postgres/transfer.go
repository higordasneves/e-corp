package postgres

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type transfer struct {
	dbPool *pgxpool.Pool
}

func NewTransferRepository(dbPool *pgxpool.Pool) repository.TransferRepo {
	return &transfer{dbPool: dbPool}
}

func (tRepo transfer) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	var db pgxtype.Querier
	db = tRepo.dbPool

	if tx := ctx.Value("dbConnection"); tx != nil {
		db = tx.(*pgxpool.Tx)
	}

	_, err := db.Exec(ctx, `INSERT INTO transfers 
		(id, account_origin_id, account_destination_id, amount, created_at)
		 VALUES ($1, $2, $3, $4, $5)`, transfer.ID.String(), transfer.AccountOriginID.String(), transfer.AccountDestinationID.String(), int(transfer.Amount), transfer.CreatedAt)

	if err != nil {
		return fmt.Errorf("unexpected sql error occurred while creating transfer: %s", err)
	}
	return nil
}

func (tRepo transfer) Transfer(ctx context.Context, ctxChan chan context.Context, errChan chan error) error {
	return PerformTransaction(ctx, ctxChan, tRepo.dbPool, errChan)
}

func (tRepo transfer) GetTransfers(ctx context.Context, cpf string) {

}
