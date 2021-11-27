package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/sirupsen/logrus"
)

type TransferUseCase interface {
	Transfer(ctx context.Context, accOriID vos.UUID, accDestID vos.UUID, amount vos.Currency)
	GetTransfers(ctx context.Context, cpf string)
}

type transferUseCase struct {
	accountRepo  repository.AccountRepo
	transferRepo repository.TransferRepo
	log          *logrus.Logger
}

func NewTransferUseCase(accountRepo repository.AccountRepo, transferRepo repository.TransferRepo, log *logrus.Logger) TransferUseCase {
	return &transferUseCase{accountRepo: accountRepo, transferRepo: transferRepo, log: log}
}
