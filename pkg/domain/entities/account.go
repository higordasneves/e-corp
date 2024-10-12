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

// AccountOutput represents information from a bank account that should be returned
type AccountOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
