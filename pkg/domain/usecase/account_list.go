package usecase

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

type ListAccountsInput struct {
	// IDs of the accounts.
	IDs []uuid.UUID
	// LastFetchedID represents the id of the last account listed in the previous page (cursor).
	LastFetchedID uuid.UUID
	// PageSize is the limit (quantity) of items that can be listed.
	PageSize int
}

type ListAccountsOutput struct {
	Accounts []entities.Account
	// NextPage is the cursor for filter the next page of accounts ad is used to create a pagination token.
	NextPage *ListAccountsInput
}

// ListAccounts Lists accounts by filtering the IDs provided in the input.
func (accUseCase AccountUseCase) ListAccounts(ctx context.Context, input ListAccountsInput) (ListAccountsOutput, error) {
	output, err := accUseCase.R.ListAccounts(ctx, input)
	if err != nil {
		return ListAccountsOutput{}, fmt.Errorf("listing accounts from db: %w", err)
	}

	return output, nil
}
