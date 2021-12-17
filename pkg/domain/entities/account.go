package entities

import (
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
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

//Account represents a banking account
type Account struct {
	ID        vos.UUID
	Name      string
	CPF       vos.CPF
	Secret    vos.Secret
	Balance   int
	CreatedAt time.Time
}

//AccountOutput represents information from a bank account that should be returned
type AccountOutput struct {
	ID        vos.UUID  `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

//GetAccOutput formats and return only pertinent information from account
func (acc *Account) GetAccOutput() *AccountOutput {
	cpf := acc.CPF.FormatOutput()
	accOutput := &AccountOutput{
		ID:        acc.ID,
		Name:      acc.Name,
		CPF:       cpf,
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt,
	}
	return accOutput
}
