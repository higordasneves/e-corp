package usecase

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type ListAccountsUCRepository interface {
	CreateAccount(ctx context.Context, acc entities.Account) error
	GetAccountByDocument(ctx context.Context, cpf vos.Document) (entities.Account, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int, error)
	ListAccounts(ctx context.Context, input ListAccountsInput) (ListAccountsOutput, error)
}

type ListAccountsUC struct {
	R ListAccountsUCRepository
}

func NewListAccountsUC(accountRepo ListAccountsUCRepository) ListAccountsUC {
	return ListAccountsUC{R: accountRepo}
}

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
func (accUseCase ListAccountsUC) ListAccounts(ctx context.Context, input ListAccountsInput) (ListAccountsOutput, error) {
	output, err := accUseCase.R.ListAccounts(ctx, input)
	if err != nil {
		return ListAccountsOutput{}, fmt.Errorf("listing accounts from db: %w", err)
	}

	return output, nil
}
