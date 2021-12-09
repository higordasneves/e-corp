package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/sirupsen/logrus"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input *AccountInput) (*entities.AccountOutput, error)
	FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
}

type accountUseCase struct {
	accountRepo repository.AccountRepo
	log         *logrus.Logger
}

func NewAccountUseCase(accountRepo repository.AccountRepo, log *logrus.Logger) AccountUseCase {
	return &accountUseCase{accountRepo: accountRepo, log: log}
}
