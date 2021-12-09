package entities

import (
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrAccNotFound      = errors.New("account not found")
	ErrAccAlreadyExists = errors.New("account already exists")
	ErrEmptyInput       = errors.New("the name, document and password fields are required")
	ErrSmallSecret      = errors.New("the password must be at least 8 characters long")
	ErrInvalidPass      = errors.New("invalid password")
)

//Account represents a banking account
type Account struct {
	ID        vos.UUID
	Name      string
	CPF       vos.CPF
	Secret    string
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

//GetHashSecret returns hash of password
func (acc *Account) GetHashSecret() error {
	hashSecret, err := bcrypt.GenerateFromPassword([]byte(acc.Secret), 10)
	if err != nil {
		return err
	}
	acc.Secret = string(hashSecret)
	return nil
}

func (acc *Account) CompareHashSecret(secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(acc.Secret), []byte(secret))
	if err != nil {
		return err
	}
	return nil
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
