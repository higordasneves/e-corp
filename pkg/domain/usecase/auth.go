package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type AuthUseCase interface {
	Login(ctx context.Context, input *LoginInput) (*Token, error)
}

type authUseCase struct {
	accountRepo repository.AccountRepo
	log         *logrus.Logger
	duration    time.Duration
	secretKey   string
}

func NewAuthUseCase(accountRepo repository.AccountRepo, log *logrus.Logger, cfgAuth *config.AuthConfig) AuthUseCase {
	return &authUseCase{accountRepo: accountRepo, log: log, duration: cfgAuth.Duration, secretKey: cfgAuth.SecretKey}
}
