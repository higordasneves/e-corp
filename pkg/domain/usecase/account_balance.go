package usecase

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

// GetBalance returns a bank account balance
func (accUseCase AccountUseCase) GetBalance(ctx context.Context, id uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	balance, err := accUseCase.accountRepo.GetBalance(ctx, id)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
