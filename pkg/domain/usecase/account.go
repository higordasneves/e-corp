package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/repository"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input AccountInput) (*models.Account, error)
}

type accountUseCase struct {
	accountRepo repository.AccountRepo
}

func NewAccountUseCase(accountRepo repository.AccountRepo) AccountUseCase {
	return &accountUseCase{accountRepo: accountRepo}
}
