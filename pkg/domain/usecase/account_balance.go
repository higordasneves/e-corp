package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

func (accUseCase accountUseCase) GetBalance(ctx context.Context, id vos.AccountID) (*vos.Currency, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	balance, err := accUseCase.accountRepo.GetBalance(ctx, id)

	if err != nil {
		return nil, err
	}
	return balance, nil
}
