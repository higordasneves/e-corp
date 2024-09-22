package usecase

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

func (tUseCase TransferUseCase) FetchTransfers(ctx context.Context, accID string) ([]entities.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	_, err := tUseCase.repo.GetBalance(ctx, uuid.FromStringOrNil(accID))
	if err != nil {
		return nil, err
	}

	return tUseCase.repo.FetchTransfers(ctx, uuid.FromStringOrNil(accID))
}
