package usecase

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestAuthUseCase_Validate(t *testing.T) {
	t.Parallel()

	accounts := []entities.Account{
		{
			CPF:    "44455566677",
			Secret: "123456",
		},
		{
			CPF:    "44455566678",
			Secret: "654321",
		},
	}

	for i, account := range accounts {
		secret, err := vos.GetHashSecret(string(account.Secret))
		if err != nil {
			logrus.Error(err)
		}
		accounts[i].Secret = secret
	}

	tests := []struct {
		name        string
		login       LoginInput
		expectedErr error
		setup       func(ctx context.Context, input LoginInput, authUC AuthUseCase) Token
	}{
		{
			name: "with success",
			login: LoginInput{
				CPF:    "44455566677",
				Secret: "123456",
			},
			expectedErr: nil,
			setup: func(ctx context.Context, input LoginInput, authUC AuthUseCase) Token {
				token, err := authUC.Login(ctx, &input)
				if err != nil {
					logrus.Fatal("unexpected login error")
				}
				return *token
			},
		},
		{
			name: "invalid token",
			login: LoginInput{
				CPF:    "44455566678",
				Secret: "654321",
			},
			expectedErr: ErrTokenInvalid,
			setup: func(ctx context.Context, input LoginInput, authUC AuthUseCase) Token {
				return "invalid_token"
			},
		},
	}

	cfg := config.Config{}
	cfg.LoadEnv()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accRepo := repomock.NewAccountRepo(accounts, tt.expectedErr)
			authUC := NewAuthUseCase(accRepo, &cfg.Auth)
			token := tt.setup(context.Background(), tt.login, authUC)

			// execute
			_, err := authUC.ValidateToken(string(token))

			// assert
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("got error: %v, want error: %v", err, tt.expectedErr)
			}
		})
	}
}
