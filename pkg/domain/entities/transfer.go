package entities

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrOriginAccID               = errors.New("invalid origin account ID")
	ErrDestAccID                 = errors.New("invalid destination account ID")
	ErrSelfTransfer              = errors.New("the destination account must be different from the origin account")
	ErrTransferAmount            = errors.New("invalid transfer amount, the amount must be greater than 0")
	ErrTransferInsufficientFunds = errors.New("insufficient funds")
)

// Transfer represents a banking transfer
type Transfer struct {
	ID                   uuid.UUID `json:"id"`
	AccountOriginID      uuid.UUID `json:"account_origin_id"`
	AccountDestinationID uuid.UUID `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
