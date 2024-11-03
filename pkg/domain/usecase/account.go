package usecase

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

//go:generate moq -stub -pkg mocks -out mocks/accounts.go . AccountUseCaseUseCaseBroker

type AccountUseCaseRepository interface {
	CreateAccount(ctx context.Context, acc entities.Account) error
	GetAccountByDocument(ctx context.Context, cpf vos.Document) (entities.Account, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int, error)
	ListAccounts(ctx context.Context, input ListAccountsInput) (ListAccountsOutput, error)
}

type AccountUseCaseUseCaseBroker interface {
	NotifyAccountCreation(ctx context.Context, account entities.Account) error
}

type AccountUseCase struct {
	R AccountUseCaseRepository
	B AccountUseCaseUseCaseBroker
}

func NewAccountUseCase(accountRepo AccountUseCaseRepository, broker AccountUseCaseUseCaseBroker) AccountUseCase {
	return AccountUseCase{R: accountRepo, B: broker}
}
