package entities

import (
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

//Transfer represents a banking transfer
type Transfer struct {
	ID                   vos.UUID
	AccountOriginID      vos.UUID
	AccountDestinationID vos.UUID
	Amount               int
	CreatedAt            time.Time
}
