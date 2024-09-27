package entities

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Transfer represents a banking transfer
type Transfer struct {
	ID                   uuid.UUID `json:"id"`
	AccountOriginID      uuid.UUID `json:"account_origin_id"`
	AccountDestinationID uuid.UUID `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
