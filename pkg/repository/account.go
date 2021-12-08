package repository

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

var (
	ErrCreateAcc  = errors.New("an unexpected error occurred while creating account")
	ErrFetchAcc   = errors.New("an unexpected error occurred while fetching accounts")
	ErrGetBalance = errors.New("an unexpected error occurred while getting account balance")
	ErrGetAccount = errors.New("an unexpected error occurred")
	ErrTruncDB    = errors.New("an unexpected error occurred while deleting tables")
)

type AccountRepo interface {
	CreateAccount(context.Context, *models.Account) error
	FetchAccounts(ctx context.Context) ([]models.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (*vos.Currency, error)
	GetAccount(ctx context.Context, cpf vos.CPF) (*models.Account, error)
}
