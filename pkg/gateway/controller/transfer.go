package controller

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

//go:generate moq -stub -pkg mocks -out mocks/transfers_uc.go . TransferUseCase

type TransferUseCase interface {
	ListAccountTransfers(ctx context.Context, input usecase.ListAccountTransfersInput) (usecase.ListAccountTransfersOutput, error)
	Transfer(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error)
}

type TransferController struct {
	tUseCase TransferUseCase
	log      *logrus.Logger
}

func NewTransferController(tUseCase TransferUseCase, log *logrus.Logger) TransferController {
	return TransferController{tUseCase: tUseCase, log: log}
}
