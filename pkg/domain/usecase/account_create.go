package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountInput struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"secret"`
	Balance float64 `json:"balance"`
}

//CreateAccount validates and handles user input and creates a formatted account
func (accUseCase accountUseCase) CreateAccount(ctx context.Context, accInput AccountInput) (*models.Account, error) {
	accID := newAccID()
	account := &models.Account{ID: vos.AccountID(accID),
		Name:    accInput.Name,
		CPF:     accInput.CPF,
		Secret:  accInput.Secret,
		Balance: accInput.Balance}

	account.GetHashSecret()

	err := accUseCase.accountRepo.CreateAccount(ctx, account)

	if err != nil {
		return nil, err
	}
	return account, nil
}

// newAccID gets uuid using google lib
func newAccID() string {
	return uuid.NewString()
}
