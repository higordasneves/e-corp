package usecase

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// GetBalance returns the current balance of the account.
// the repository return domain.ErrNotFound if the account not exists.
func (accUseCase AccountUseCase) GetBalance(ctx context.Context, id uuid.UUID) (int, error) {
	balance, err := accUseCase.R.GetBalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("error getting account balance: %w", err)
	}

	return balance, nil
}
