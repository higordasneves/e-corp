package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestAuthUseCase_Login(t *testing.T) {
	t.Parallel()

	//setup
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
	}{
		{
			name: "with success",
			login: LoginInput{
				CPF:    "44455566677",
				Secret: "123456",
			},
			expectedErr: nil,
		},
		{
			name: "with success 2",
			login: LoginInput{
				CPF:    "44455566678",
				Secret: "654321",
			},
			expectedErr: nil,
		},
		{
			name: "account not found",
			login: LoginInput{
				CPF:    "44455566679",
				Secret: "secret",
			},
			expectedErr: entities.ErrAccNotFound,
		},
		{
			name: "invalid password",
			login: LoginInput{
				CPF:    "44455566678",
				Secret: "wrong_secret",
			},
			expectedErr: vos.ErrInvalidPass,
		},
		{
			name: "database error",
			login: LoginInput{
				CPF:    "44455566679",
				Secret: "secret",
			},
			expectedErr: repository.ErrUnexpected,
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

			//execute
			ctx := context.Background()
			_, err := authUC.Login(ctx, &tt.login)

			//assert
			if err != tt.expectedErr {
				t.Errorf("got error: %v, want error: %v", err, tt.expectedErr)
			}
		})
	}
}
