package usecase

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// GetBalance returns a bank account balance
func (accUseCase AccountUseCase) GetBalance(ctx context.Context, id uuid.UUID) (int, error) {
	balance, err := accUseCase.R.GetBalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("error getting account balance: %w", err)
	}

	return balance, nil
}
