package ucmock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

const balanceInit = 1000000

type AccountUseCase interface {
	CreateAccount(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error)
	FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
}

type accountUseCase struct {
	accounts  []entities.AccountOutput
	domainErr error
}

func NewAccountUseCase(accounts []entities.AccountOutput, domainErr error) AccountUseCase {
	return &accountUseCase{accounts: accounts, domainErr: domainErr}
}

func (accUseCase accountUseCase) CreateAccount(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
	if accUseCase.domainErr != nil {
		return nil, accUseCase.domainErr
	}

	cpf := input.CPF.FormatOutput()
	accOutput := &entities.AccountOutput{
		ID:        vos.NewUUID(),
		Name:      input.Name,
		CPF:       cpf,
		Balance:   balanceInit,
		CreatedAt: time.Now().Truncate(time.Second),
	}

	return accOutput, nil
}

func (accUseCase accountUseCase) FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error) {
	if accUseCase.domainErr != nil {
		return nil, accUseCase.domainErr
	}

	return accUseCase.accounts, nil
}

func (accUseCase accountUseCase) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	if accUseCase.domainErr != nil {
		return 0, accUseCase.domainErr
	}

	for _, acc := range accUseCase.accounts {
		if acc.ID == id {
			return acc.Balance, nil
		}
	}

	return 0, entities.ErrAccNotFound
}
