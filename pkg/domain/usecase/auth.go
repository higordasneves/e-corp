package usecase

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/repository"
	"time"
)

type AuthUseCase interface {
	Login(ctx context.Context, input *LoginInput) (*Token, error)
	ValidateToken(tokenString string) (*jwt.StandardClaims, error)
}

type authUseCase struct {
	accountRepo repository.AccountRepo
	duration    time.Duration
	secretKey   string
}

func NewAuthUseCase(accountRepo repository.AccountRepo, cfgAuth *config.AuthConfig) AuthUseCase {
	return &authUseCase{accountRepo: accountRepo, duration: cfgAuth.Duration, secretKey: cfgAuth.SecretKey}
}
