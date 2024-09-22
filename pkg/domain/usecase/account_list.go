package usecase

import (
	"context"
	"time"

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

// FetchAccounts calls the func to select all accounts
func (accUseCase AccountUseCase) FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	output, err := accUseCase.R.ListAccounts(ctx, ListAccountsInput{})
	if err != nil {
		return nil, err
	}
	accListOutput := make([]entities.AccountOutput, 0, len(output.Accounts))
	for _, acc := range output.Accounts {
		out := acc.GetAccOutput()
		accListOutput = append(accListOutput, *out)
	}
	return accListOutput, nil
}
