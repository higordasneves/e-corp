package models

import (
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

//Account represents a banking account
type Account struct {
	ID        vos.AccountID
	Name      string
	CPF       string
	Secret    string
	Balance   vos.Currency
	CreatedAt time.Time
}

//AccountOutput represents information from a bank account that should be returned
type AccountOutput struct {
	ID        vos.AccountID `json:"id"`
	Name      string        `json:"name"`
	CPF       string        `json:"cpf"`
	Balance   vos.Currency  `json:"balance"`
	CreatedAt time.Time     `json:"created_at"`
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

//GetAccOutput formats and return only pertinent information from account
func (acc *Account) GetAccOutput() *AccountOutput {
	acc.cpfMask()
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

//cpfMask formats CPF of account owner to pattern "xxx-xxx-xxx-xx"
func (acc *Account) cpfMask() {
	cpfModel, err := regexp.Compile(`^([\d]{3})([\d]{3})([\d]{3})([\d]{2})$`)
	if err == nil {
		acc.CPF = cpfModel.ReplaceAllString(acc.CPF, "$1.$2.$3-$4")
	}
}
