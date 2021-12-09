package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

func (accUseCase accountUseCase) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	balance, err := accUseCase.accountRepo.GetBalance(ctx, id)

	if err != nil {
		return 0, err
	}

	return balance, nil
}
