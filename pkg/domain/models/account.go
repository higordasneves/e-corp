package models

import (
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//Account represents a banking account
type Account struct {
	ID        vos.UUID
	Name      string
	CPF       vos.CPF
	Secret    string
	Balance   vos.Currency
	CreatedAt time.Time
}

//AccountOutput represents information from a bank account that should be returned
type AccountOutput struct {
	ID        vos.UUID     `json:"id"`
	Name      string       `json:"name"`
	CPF       vos.CPF      `json:"cpf"`
	Balance   vos.Currency `json:"balance"`
	CreatedAt time.Time    `json:"created_at"`
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
	acc.CPF.FormatOutput()
	acc.Balance.ConvertFromCents()
	accOutput := &AccountOutput{
		ID:        acc.ID,
		Name:      acc.Name,
		CPF:       acc.CPF,
		Balance:   acc.Balance,
		CreatedAt: acc.CreatedAt,
	}
	return accOutput
}
