package entities

import (
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

var (
	ErrOriginAccID               = errors.New("invalid origin account ID")
	ErrDestAccID                 = errors.New("invalid destination account ID")
	ErrSelfTransfer              = errors.New("the destination account must be different from the origin account")
	ErrTransferAmount            = errors.New("invalid transfer amount, value must be greater than 0")
	ErrTransferInsufficientFunds = errors.New("insufficient funds")
)

//Transfer represents a banking transfer
type Transfer struct {
	ID                   vos.UUID
	AccountOriginID      vos.UUID
	AccountDestinationID vos.UUID
	Amount               int
	CreatedAt            time.Time
}
