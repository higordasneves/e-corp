package repomock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type TransferRepo interface {
	CreateTransfer(ctx context.Context, transfer *entities.Transfer) error
	PerformTransaction(ctx context.Context, ctxChan chan context.Context, errChan chan error) error
	FetchTransfers(ctx context.Context, id vos.UUID) ([]entities.Transfer, error)
}

type transfer struct {
	transfers []entities.Transfer
	dbError   error
}

func NewTransferRepo(transfers []entities.Transfer, dbError error) TransferRepo {
	return &transfer{transfers: transfers, dbError: dbError}
}

func (tRepo transfer) CreateTransfer(context.Context, *entities.Transfer) error {
	return tRepo.dbError
}

func (tRepo transfer) PerformTransaction(ctx context.Context, ctxChan chan context.Context, errChan chan error) error {
	ctxChan <- ctx
	return <-errChan
}

func (tRepo transfer) FetchTransfers(_ context.Context, id vos.UUID) ([]entities.Transfer, error) {
	if tRepo.dbError != nil {
		return nil, tRepo.dbError
	}

	var accTransfer []entities.Transfer
	for _, t := range tRepo.transfers {
		if t.AccountOriginID == id {
			accTransfer = append(accTransfer, t)
		}
	}

	return accTransfer, nil
}
