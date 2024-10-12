package controller

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

//go:generate moq -stub -pkg mocks -out mocks/accounts_uc.go . AccountUseCase

type AccountController struct {
	accUseCase AccountUseCase
	log        *logrus.Logger
}

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int, error)
	ListAccounts(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error)
}

func NewAccountController(accUseCase AccountUseCase, log *logrus.Logger) AccountController {
	return AccountController{accUseCase: accUseCase, log: log}
}
