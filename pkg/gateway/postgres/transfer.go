package postgres

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/sqlc"
)

// CreateTransfer inserts a transfer in the database.
func (r Repository) CreateTransfer(ctx context.Context, transfer entities.Transfer) error {
	err := sqlc.New(r.conn.GetTxOrPool(ctx)).InsertTransfer(ctx, sqlc.InsertTransferParams{
		ID:                   uuid.FromStringOrNil(transfer.ID.String()),
		AccountOriginID:      uuid.FromStringOrNil(transfer.AccountOriginID.String()),
		AccountDestinationID: uuid.FromStringOrNil(transfer.AccountDestinationID.String()),
		Amount:               int64(transfer.Amount),
		CreatedAt:            transfer.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("inserting transfer with id %s: %w", transfer.ID.String(), err)
	}

	return nil
}

// ListAccountTransfers lists all transfers made or received by an account in descending order.
func (r Repository) ListAccountTransfers(ctx context.Context, accountID uuid.UUID) ([]entities.Transfer, error) {
	rows, err := sqlc.New(r.conn.GetTxOrPool(ctx)).ListAccountTransfers(ctx, uuid.FromStringOrNil(accountID.String()))
	if err != nil {
		return nil, fmt.Errorf("listing transfer for account %s: %w", accountID.String(), err)
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
