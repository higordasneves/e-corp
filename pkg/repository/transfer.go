package repository

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type TransferRepo interface {
	CreateTransfer(ctx context.Context, transfer *entities.Transfer) error
	Transfer(ctx context.Context, ctxChan chan context.Context, errChan chan error) error
	FetchTransfers(ctx context.Context, id vos.UUID) ([]entities.Transfer, error)
}
