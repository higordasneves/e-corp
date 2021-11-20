package models

import (
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AccountID string

//Account represents a banking account
type Account struct {
	ID        vos.AccountID
	Name      string
	CPF       string
	Secret    string
	Balance   float64
	CreatedAt time.Time
}

func (acc *Account) GetHashSecret() error {
	hashSecret, err := bcrypt.GenerateFromPassword([]byte(acc.Secret), 10)
	if err != nil {
		return err
	}
	acc.Secret = string(hashSecret)
	return nil
}
