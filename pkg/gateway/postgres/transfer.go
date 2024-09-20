package postgres

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/sqlc"
)

func (r Repository) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	var db Querier
	db = r.dbPool

	if tx := ctx.Value("dbConnection"); tx != nil {
		db = tx.(*pgxpool.Tx)
	}

	err := sqlc.New(db).InsertTransfer(ctx, sqlc.InsertTransferParams{
		ID:                   uuid.FromStringOrNil(transfer.ID.String()),
		AccountOriginID:      uuid.FromStringOrNil(transfer.AccountOriginID.String()),
		AccountDestinationID: uuid.FromStringOrNil(transfer.AccountDestinationID.String()),
		Amount:               int64(transfer.Amount),
		CreatedAt:            transfer.CreatedAt,
	})
	if err != nil {
		return domain.NewDBError(domain.QueryRefCreateTransfer, err, domain.ErrUnexpected)
	}

	return nil
}

func (r Repository) PerformTransaction(ctx context.Context, ctxChan chan context.Context, errChan chan error) error {
	return PerformTransaction(ctx, ctxChan, r.dbPool, errChan)
}

func (r Repository) FetchTransfers(ctx context.Context, id vos.UUID) ([]entities.Transfer, error) {
	rows, err := sqlc.New(r.dbPool).ListAccountSentTransfers(ctx, uuid.FromStringOrNil(id.String()))
	if err != nil {
		return nil, domain.NewDBError(domain.QueryRefGetTransfers, err, domain.ErrUnexpected)
	}

	transferList := make([]entities.Transfer, 0, len(rows))
	for _, row := range rows {
		transferList = append(transferList, parseSqlcTransfer(row))
	}

	return transferList, nil
}

func parseSqlcTransfer(t sqlc.Transfer) entities.Transfer {
	return entities.Transfer{
		ID:                   vos.UUID(t.ID.String()),
		AccountOriginID:      vos.UUID(t.AccountOriginID.String()),
		AccountDestinationID: vos.UUID(t.AccountDestinationID.String()),
		Amount:               int(t.Amount),
		CreatedAt:            t.CreatedAt,
	}
}
