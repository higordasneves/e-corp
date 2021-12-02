package repomock

import (
	"context"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

type AccountRepo interface {
	CreateAccount(context.Context, *models.Account) error
	FetchAccounts(ctx context.Context) ([]models.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (*vos.Currency, error)
	GetAccount(ctx context.Context, cpf string) (*models.Account, error)
}

type account struct {
	accounts []models.Account
	err      error
}

func NewAccountRepo(accounts []models.Account, err error) AccountRepo {
	return &account{accounts, err}
}

func (accRepo account) CreateAccount(context.Context, *models.Account) error {
	return accRepo.err
}

func (accRepo account) FetchAccounts(context.Context) ([]models.Account, error) {
	if accRepo.err != nil {
		return nil, accRepo.err
	}

	accountsList := make([]models.Account, 0, len(accRepo.accounts))
	for _, acc := range accRepo.accounts {
		accountOutput := models.Account{
			ID:        acc.ID,
			Name:      acc.Name,
			CPF:       acc.CPF,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
		}
		accountsList = append(accountsList, accountOutput)
	}
	return accountsList, nil
}

func (accRepo account) GetBalance(ctx context.Context, id vos.UUID) (*vos.Currency, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if accRepo.err != nil {
		return nil, accRepo.err
	}

	for _, acc := range accRepo.accounts {
		if id == acc.ID {
			return &acc.Balance, nil
		}
	}
	return nil, domainerr.ErrAccNotFound
}

func (accRepo account) GetAccount(context.Context, string) (*models.Account, error) {
	panic("implement me")
}
