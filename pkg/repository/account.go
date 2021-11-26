package repository

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

var (
	ErrCreateAcc  = errors.New("an unexpected error has occurred while creating account")
	ErrFetchAcc   = errors.New("an unexpected error occurred while fetching accounts")
	ErrGetBalance = errors.New("an unexpected error occurred while getting account balance")
)

type AccountRepo interface {
	CreateAccount(context.Context, *models.Account) error
	FetchAccounts(ctx context.Context) ([]models.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.AccountID) (*vos.Currency, error)
}
