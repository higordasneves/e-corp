package controller

import (
	"context"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

//go:generate moq -stub -pkg mocks -out mocks/transfers_uc.go . TransferUseCase

type TransferUseCase interface {
	ListAccountTransfers(ctx context.Context, input usecase.ListAccountTransfersInput) (usecase.ListAccountTransfersOutput, error)
	Transfer(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error)
}

type TransferController struct {
	tUseCase TransferUseCase
}

func NewTransferController(tUseCase TransferUseCase) TransferController {
	return TransferController{tUseCase: tUseCase}
}
