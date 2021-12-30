package ucmock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountUseCase struct {
	Create        func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error)
	Fetch         func(ctx context.Context) ([]entities.AccountOutput, error)
	GetAccBalance func(ctx context.Context, id vos.UUID) (int, error)
}

func (accUseCase AccountUseCase) CreateAccount(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
	return accUseCase.Create(ctx, input)
}

func (accUseCase AccountUseCase) FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error) {
	return accUseCase.Fetch(ctx)
}

func (accUseCase AccountUseCase) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	return accUseCase.GetAccBalance(ctx, id)
}
