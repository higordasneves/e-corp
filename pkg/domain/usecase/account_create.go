package usecase

import (
	"github.com/google/uuid"
	"github.com/higordasneves/e-corp/pkg/domain/models"
)

type AccountInput struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"secret"`
	Balance float64 `json:"balance"`
}

func (a accountUseCase) CreateAccount(accInput AccountInput) (*models.Account, error) {
	accID := newAccID()
	account := &models.Account{ID: accID,
		Name:    accInput.Name,
		CPF:     accInput.CPF,
		Secret:  accInput.Secret,
		Balance: accInput.Balance}
	account.GetHashSecret()
	return account, nil
}

func newAccID() string {
	return uuid.NewString()
}
