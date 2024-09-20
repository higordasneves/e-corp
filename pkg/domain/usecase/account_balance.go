package usecase

import (
	"context"
	"time"

	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

// GetBalance returns a bank account balance
func (accUseCase AccountUseCase) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := vos.IsValidUUID(id.String())
	if err != nil {
		return 0, err
	}

	balance, err := accUseCase.accountRepo.GetBalance(ctx, id)

	if err != nil {
		return 0, err
	}

	return balance, nil
}
