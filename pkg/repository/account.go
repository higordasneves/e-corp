package repository

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
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
	CreateAccount(context.Context, *entities.Account) error
	FetchAccounts(ctx context.Context) ([]entities.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
	GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
}
