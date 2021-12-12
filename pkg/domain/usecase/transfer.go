package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/repository"
)

type TransferUseCase interface {
	Transfer(ctx context.Context, transferInput *TransferInput) (*entities.Transfer, error)
	GetTransfers(ctx context.Context, cpf string)
}

type transferUseCase struct {
	accountRepo  repository.AccountRepo
	transferRepo repository.TransferRepo
}

func NewTransferUseCase(accountRepo repository.AccountRepo, transferRepo repository.TransferRepo) TransferUseCase {
	return &transferUseCase{accountRepo: accountRepo, transferRepo: transferRepo}
}
