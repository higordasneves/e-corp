package repository

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type TransferRepo interface {
	Transfer(ctx context.Context, accOriID vos.UUID, accDestID vos.UUID, amount vos.Currency)
	GetTransfers(ctx context.Context, cpf string)
}
