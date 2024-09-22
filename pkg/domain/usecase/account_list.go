package usecase

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

type ListAccountsInput struct {
	IDs           []uuid.UUID
	LastFetchedID uuid.UUID
	PageSize      int
}

type ListAccountsOutput struct {
	Accounts []entities.Account
	NextPage *ListAccountsInput
}

// FetchAccounts calls the func to select all accounts
func (accUseCase AccountUseCase) FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	output, err := accUseCase.accountRepo.ListAccounts(ctx, ListAccountsInput{})
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
