package usecase

import (
	"time"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
)

type AuthUseCase struct {
	accountRepo AccountUseCaseRepository
	duration    time.Duration
	secretKey   string
}

func NewAuthUseCase(accountRepo AccountUseCaseRepository, cfgAuth *config.AuthConfig) AuthUseCase {
	return AuthUseCase{accountRepo: accountRepo, duration: cfgAuth.Duration, secretKey: cfgAuth.SecretKey}
}
