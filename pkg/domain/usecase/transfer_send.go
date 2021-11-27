package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func (tUseCase transferUseCase) Transfer(ctx context.Context, accOriID vos.UUID, accDestID vos.UUID, amount vos.Currency) {

}
