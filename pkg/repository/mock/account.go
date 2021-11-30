package repomock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountRepo interface {
	CreateAccount(context.Context, *models.Account) error
	FetchAccounts(ctx context.Context) ([]models.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (*vos.Currency, error)
	GetAccount(ctx context.Context, cpf string) (*models.Account, error)
}

type account struct {
	acc *models.Account
	err error
}

func NewAccountRepo(accRepo *models.Account, err error) AccountRepo {
	return &account{accRepo, err}
}

func (accRepo account) CreateAccount(context.Context, *models.Account) error {
	return accRepo.err
}

func (accRepo account) FetchAccounts(context.Context) ([]models.Account, error) {
	panic("implement me")
}

func (accRepo account) GetBalance(context.Context, vos.UUID) (*vos.Currency, error) {
	panic("implement me")
}

func (accRepo account) GetAccount(context.Context, string) (*models.Account, error) {
	panic("implement me")
}
