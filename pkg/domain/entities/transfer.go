package entities

import (
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

var (
	//ErrBadTransferRequest generic error for bad requests
	ErrBadTransferRequest = errors.New("bad request")
	ErrOriginAccID        = errors.New("invalid origin account ID")
	ErrDestAccID          = errors.New("invalid destination account ID")
	ErrTransferAmount     = errors.New("invalid transfer amount, value must be greater than 0")
)

//Transfer represents a banking transfer
type Transfer struct {
	ID                   vos.UUID
	AccountOriginID      vos.UUID
	AccountDestinationID vos.UUID
	Amount               int
	CreatedAt            time.Time
}
