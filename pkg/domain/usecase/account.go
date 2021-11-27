package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/sirupsen/logrus"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input AccountInput) (*models.Account, error)
	FetchAccounts(ctx context.Context) ([]models.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.UUID) (*vos.Currency, error)
}

type accountUseCase struct {
	accountRepo repository.AccountRepo
	log         *logrus.Logger
}

func NewAccountUseCase(accountRepo repository.AccountRepo, log *logrus.Logger) AccountUseCase {
	return &accountUseCase{accountRepo: accountRepo, log: log}
}
