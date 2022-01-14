package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
)

//go:generate moq -skip-ensure -stub -out mock/account.go -pkg ucmock ./../../domain/usecase AccountUseCase:AccountUseCase

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input *AccountInput) (*entities.AccountOutput, error)
	FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
}

type accountUseCase struct {
	accountRepo repository.AccountRepo
}

func NewAccountUseCase(accountRepo repository.AccountRepo) AccountUseCase {
	return &accountUseCase{accountRepo: accountRepo}
}
