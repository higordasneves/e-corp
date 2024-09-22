package usecase

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type TransferUseCaseRepository interface {
	GetBalance(ctx context.Context, id uuid.UUID) (int, error)
	GetAccountByDocument(ctx context.Context, cpf vos.CPF) (entities.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, transactionAmount int) error

	CreateTransfer(ctx context.Context, transfer *entities.Transfer) error
	FetchTransfers(ctx context.Context, id uuid.UUID) ([]entities.Transfer, error)

	BeginTX(ctx context.Context) (context.Context, error)
	CommitTX(ctx context.Context) error
	RollbackTX(ctx context.Context) error
}

type TransferUseCase struct {
	repo TransferUseCaseRepository
}

func NewTransferUseCase(r TransferUseCaseRepository) TransferUseCase {
	return TransferUseCase{repo: r}
}
