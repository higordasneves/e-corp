package repository

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

type TransferRepo interface {
	CreateTransfer(ctx context.Context, transfer *entities.Transfer) error
	GetTransfers(ctx context.Context, cpf string)
	Transfer(ctx context.Context, ctxChan chan context.Context, errChan chan error) error
}
