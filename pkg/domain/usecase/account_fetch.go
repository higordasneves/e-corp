package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
)

//FetchAccounts calls the func to select all accounts
func (accUseCase *accountUseCase) FetchAccounts(ctx context.Context) ([]models.AccountOutput, error) {
	accList, err := accUseCase.accountRepo.FetchAccounts(ctx)
	if err != nil {
		return nil, err
	}
	return accList, nil
}
