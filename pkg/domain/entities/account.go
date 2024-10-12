package entities

import (
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

// Account represents a banking account
type Account struct {
	ID        uuid.UUID
	Name      string
	Document  vos.Document
	Secret    vos.Secret
	Balance   int
	CreatedAt time.Time
}
