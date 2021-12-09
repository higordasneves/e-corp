package repomock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

type AccountRepo interface {
	CreateAccount(context.Context, *entities.Account) error
	FetchAccounts(ctx context.Context) ([]entities.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
	GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
}

type account struct {
	accounts []entities.Account
	err      error
}

func NewAccountRepo(accounts []entities.Account, err error) AccountRepo {
	return &account{accounts, err}
}

func (accRepo account) CreateAccount(context.Context, *entities.Account) error {
	return accRepo.err
}

func (accRepo account) FetchAccounts(context.Context) ([]entities.Account, error) {
	if accRepo.err != nil {
		return nil, accRepo.err
	}

	accountsList := make([]entities.Account, 0, len(accRepo.accounts))
	for _, acc := range accRepo.accounts {
		accountOutput := entities.Account{
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

func (accRepo account) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if accRepo.err != nil {
		return 0, accRepo.err
	}

	for _, acc := range accRepo.accounts {
		if id == acc.ID {
			return acc.Balance, nil
		}
	}
	return 0, entities.ErrAccNotFound
}

func (accRepo account) GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	panic("implement me")
}
