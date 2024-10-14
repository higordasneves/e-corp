package usecase

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

type ListAccountTransfersInput struct {
	AccountID uuid.UUID
}

type ListAccountTransfersOutput struct {
	Transfers []entities.Transfer
}

// ListAccountTransfers lists all the transfers sent or received by the account in desc order.
// Returns domain.ErrNotFound if the account not exists.
func (tUseCase TransferUseCase) ListAccountTransfers(ctx context.Context, input ListAccountTransfersInput) (ListAccountTransfersOutput, error) {
	// Just checking if the account exists. the repository returns domain.ErrNotFound if not exits.
	if _, err := tUseCase.R.GetBalance(ctx, input.AccountID); err != nil {
		return ListAccountTransfersOutput{}, fmt.Errorf("getting balance: %w", err)
	}

	output, err := tUseCase.R.ListAccountTransfers(ctx, input.AccountID)
	if err != nil {
		return ListAccountTransfersOutput{}, fmt.Errorf("listing transfers: %w", err)
	}

	return ListAccountTransfersOutput{output}, nil
}
