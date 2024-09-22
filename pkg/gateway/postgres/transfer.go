package postgres

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/sqlc"
)

func (r Repository) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	err := sqlc.New(r.conn.GetTxOrPool(ctx)).InsertTransfer(ctx, sqlc.InsertTransferParams{
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

func (r Repository) FetchTransfers(ctx context.Context, id uuid.UUID) ([]entities.Transfer, error) {
	rows, err := sqlc.New(r.conn.GetTxOrPool(ctx)).ListAccountSentTransfers(ctx, uuid.FromStringOrNil(id.String()))
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
		ID:                   t.ID,
		AccountOriginID:      t.AccountOriginID,
		AccountDestinationID: t.AccountDestinationID,
		Amount:               int(t.Amount),
		CreatedAt:            t.CreatedAt,
	}
}
