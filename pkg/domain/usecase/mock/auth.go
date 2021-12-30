package ucmock

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

type AuthUseCase struct {
	AuthLogin func(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error)
	Validate  func(tokenString string) (*jwt.StandardClaims, error)
}

func (authUC AuthUseCase) Login(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error) {
	return authUC.AuthLogin(ctx, input)
}

func (authUC AuthUseCase) ValidateToken(tokenString string) (*jwt.StandardClaims, error) {
	return authUC.Validate(tokenString)
}
