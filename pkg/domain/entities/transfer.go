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
	//ErrZeroRowsAffectedCreateTransfer occurs when zero rows affected in create transfer query
	ErrZeroRowsAffectedCreateTransfer = errors.New("zero rows affected in create transfer query")
)

//Transfer represents a banking transfer
type Transfer struct {
	ID                   vos.UUID  `json:"id"`
	AccountOriginID      vos.UUID  `json:"account_origin_id"`
	AccountDestinationID vos.UUID  `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
