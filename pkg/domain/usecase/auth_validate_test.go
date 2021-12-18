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
	ctx := context.Background()

	cfg := config.Config{}
	cfg.LoadEnv()

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
		funcToken   func(input LoginInput, authUC AuthUseCase) Token
	}{
		{
			name: "with success",
			login: LoginInput{
				CPF:    "44455566677",
				Secret: "123456",
			},
			expectedErr: nil,
			funcToken: func(input LoginInput, authUC AuthUseCase) Token {
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
			funcToken: func(input LoginInput, authUC AuthUseCase) Token {
				return "invalid_token"
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accRepo := repomock.NewAccountRepo(accounts, test.expectedErr)
			authUC := NewAuthUseCase(accRepo, &cfg.Auth)
			_, err := authUC.Login(ctx, &test.login)

			token := test.funcToken(test.login, authUC)

			_, err = authUC.ValidateToken(string(token))

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("got error: %v, want error: %v", err, test.expectedErr)
			}
		})
	}
}
