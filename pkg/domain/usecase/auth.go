package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/sirupsen/logrus"
)

type AuthUseCase interface {
	Login(ctx context.Context, input *LoginInput) (*models.Account, error)
}

type authUseCase struct {
	accountRepo repository.AccountRepo
	log         *logrus.Logger
}

func NewAuthUseCase(accountRepo repository.AccountRepo, log *logrus.Logger) AuthUseCase {
	return &authUseCase{accountRepo: accountRepo, log: log}
}
