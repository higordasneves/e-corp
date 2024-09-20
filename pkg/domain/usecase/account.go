package usecase

import (
	"context"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountUseCaseRepository interface {
	CreateAccount(ctx context.Context, acc *entities.Account) error
	GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
	FetchAccounts(ctx context.Context) ([]entities.Account, error)
}

type AccountUseCase struct {
	accountRepo AccountUseCaseRepository
}

func NewAccountUseCase(accountRepo AccountUseCaseRepository) AccountUseCase {
	return AccountUseCase{accountRepo: accountRepo}
}
