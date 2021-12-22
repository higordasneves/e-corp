package ucmock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error)
	FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
}

type accountUseCase struct {
	accounts  []entities.Account
	domainErr error
}

func NewAccountUseCase(accounts []entities.Account, domainErr error) AccountUseCase {
	return &accountUseCase{accounts: accounts, domainErr: domainErr}
}

func (a accountUseCase) CreateAccount(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
	panic("implement me")
}

func (a accountUseCase) FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error) {
	panic("implement me")
}

func (a accountUseCase) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	panic("implement me")
}
