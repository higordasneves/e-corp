package repomock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"time"
)

type AccountRepo interface {
	CreateAccount(context.Context, *entities.Account) error
	FetchAccounts(ctx context.Context) ([]entities.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
	GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
	UpdateBalance(ctx context.Context, id vos.UUID, transactionAmount int) error
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
	if accRepo.err == repository.ErrUnexpected {
		return nil, repository.ErrUnexpected
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
	if accRepo.err == repository.ErrUnexpected {
		return 0, repository.ErrUnexpected
	}

	for _, acc := range accRepo.accounts {
		if id == acc.ID {
			return acc.Balance, nil
		}
	}
	return 0, entities.ErrAccNotFound
}

func (accRepo account) GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if accRepo.err == repository.ErrUnexpected {
		return nil, repository.ErrUnexpected
	}

	for _, acc := range accRepo.accounts {
		if cpf == acc.CPF {
			return &acc, nil
		}
	}
	return nil, entities.ErrAccNotFound
}

func (accRepo account) UpdateBalance(ctx context.Context, id vos.UUID, transactionAmount int) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if accRepo.err != nil {
		return accRepo.err
	}

	for _, acc := range accRepo.accounts {
		if id == acc.ID {
			acc.Balance += transactionAmount
			return nil
		}
	}
	return entities.ErrZeroRowsAffectedUpdateBalance
}
