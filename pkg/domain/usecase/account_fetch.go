package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"time"
)

//FetchAccounts calls the func to select all accounts
func (accUseCase *accountUseCase) FetchAccounts(ctx context.Context) ([]entities.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	accList, err := accUseCase.accountRepo.FetchAccounts(ctx)
	if err != nil {
		return nil, err
	}
	accListOutput := make([]entities.AccountOutput, 0, len(accList))
	for _, acc := range accList {
		out := acc.GetAccOutput()
		accListOutput = append(accListOutput, *out)
	}
	return accListOutput, nil
}
