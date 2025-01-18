package usecase

import (
	"context"
	"fmt"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"

	"github.com/gofrs/uuid/v5"
)

type GetAccountBalanceUCRepository interface {
	CreateAccount(ctx context.Context, acc entities.Account) error
	GetAccountByDocument(ctx context.Context, cpf vos.Document) (entities.Account, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int, error)
	ListAccounts(ctx context.Context, input ListAccountsInput) (ListAccountsOutput, error)
}

type GetAccountBalanceUC struct {
	R GetAccountBalanceUCRepository
}

func NewGetAccountBalanceUC(accountRepo GetAccountBalanceUCRepository) GetAccountBalanceUC {
	return GetAccountBalanceUC{R: accountRepo}
}

// GetBalance returns the current balance of the account.
// the repository return domain.ErrNotFound if the account not exists.
func (accUseCase GetAccountBalanceUC) GetBalance(ctx context.Context, id uuid.UUID) (int, error) {
	balance, err := accUseCase.R.GetBalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("error getting account balance: %w", err)
	}

	return balance, nil
}
