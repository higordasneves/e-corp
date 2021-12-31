package ucmock

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

type TransferUseCase struct {
	Send  func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error)
	Fetch func(ctx context.Context, id string) ([]entities.Transfer, error)
}

func (tUseCase TransferUseCase) Transfer(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
	return tUseCase.Send(ctx, transferInput)
}

func (tUseCase TransferUseCase) FetchTransfers(ctx context.Context, id string) ([]entities.Transfer, error) {
	return tUseCase.Fetch(ctx, id)
}
