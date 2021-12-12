package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"strings"
	"time"
)

func (tUseCase transferUseCase) FetchTransfers(ctx context.Context, id string) ([]entities.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	id = strings.TrimSpace(id)
	err := vos.IsValidUUID(id)
	if err != nil {
		return nil, err
	}
	accID := vos.UUID(id)

	_, err = tUseCase.accountRepo.GetBalance(ctx, accID)
	if err != nil {
		return nil, err
	}

	return tUseCase.transferRepo.FetchTransfers(ctx, accID)
}
