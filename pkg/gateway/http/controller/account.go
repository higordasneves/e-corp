package controller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

//go:generate moq -stub -pkg mocks -out mocks/accounts_uc.go . AccountUseCase

type AccountController struct {
	accUseCase AccountUseCase
	log        *logrus.Logger
}

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
	FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error)
}

func NewAccountController(accUseCase AccountUseCase, log *logrus.Logger) AccountController {
	return AccountController{accUseCase: accUseCase, log: log}
}
