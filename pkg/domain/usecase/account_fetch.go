package usecase

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
)

var (
	ErrFetchAccounts = errors.New("error fetching bank accounts")
)

//FetchAccounts calls the func to select all accounts
func (accUseCase *accountUseCase) FetchAccounts(ctx context.Context) ([]models.AccountOutput, error) {
	accList, err := accUseCase.accountRepo.FetchAccounts(ctx)
	if err != nil {
		accUseCase.log.WithError(err).Print(ErrFetchAccounts)
		return nil, ErrFetchAccounts
	}
	return accList, nil
}
