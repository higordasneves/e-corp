package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func (accUseCase accountUseCase) GetBalance(ctx context.Context, id vos.AccountID) (*vos.Currency, error) {
	balance, err := accUseCase.accountRepo.GetBalance(ctx, id)

	if err != nil {
		return nil, err
	}
	return balance, nil
}
