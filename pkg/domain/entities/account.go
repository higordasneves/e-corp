package entities

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

var (
	//ErrAccNotFound occurs when trying to obtain information from a non-existent account
	ErrAccNotFound = errors.New("account not found")
	//ErrAccAlreadyExists occurs when trying to create an account that already exists
	ErrAccAlreadyExists = errors.New("account already exists")
	//ErrEmptyInput occurs when fields required to create an account aren't filled
	ErrEmptyInput = errors.New("the name, document and password fields are required")
	//ErrZeroRowsAffectedUpdateBalance occurs when zero rows affected in update balance query
	ErrZeroRowsAffectedUpdateBalance = errors.New("zero rows affected in update balance query")
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
