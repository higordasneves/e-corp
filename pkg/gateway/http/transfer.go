package http

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

//go:generate moq -stub -pkg mocks -out mocks/transfers_uc.go . TransferUseCase

type TransferUseCase interface {
	Transfer(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error)
	FetchTransfers(ctx context.Context, id string) ([]entities.Transfer, error)
}

type TransferController struct {
	tUseCase TransferUseCase
	log      *logrus.Logger
}

func NewTransferController(tUseCase TransferUseCase, log *logrus.Logger) TransferController {
	return TransferController{tUseCase: tUseCase, log: log}
}
