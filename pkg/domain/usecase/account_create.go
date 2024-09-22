package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

const balanceInit = 1000000

// AccountInput represents information necessary to create a bank account
type AccountInput struct {
	Name   string  `json:"name"`
	CPF    vos.CPF `json:"cpf"`
	Secret string  `json:"secret"`
}

// CreateAccount validates and handles user input and creates a formatted account,
// then calls the function to insert the account into the database
func (accUseCase AccountUseCase) CreateAccount(ctx context.Context, accInput *AccountInput) (*entities.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := accInput.ValidateAccountInput()
	if err != nil {
		return nil, err
	}

	hashSecret, err := vos.GetHashSecret(accInput.Secret)
	if err != nil {
		return nil, err
	}

	account := entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      accInput.Name,
		CPF:       accInput.CPF,
		Secret:    hashSecret,
		Balance:   balanceInit,
		CreatedAt: time.Now().Truncate(time.Second),
	}

	err = accUseCase.R.CreateAccount(ctx, account)

	if err != nil {
		return nil, err
	}

	return account.GetAccOutput(), nil
}

// ValidateAccountInput validates account input and returns if occurred an error
func (accInput *AccountInput) ValidateAccountInput() error {
	accInput.removeBlankSpaces()

	err := accInput.validateInputEmpty()
	if err != nil {
		return err
	}

	err = vos.ValidateSecretLen(accInput.Secret)
	if err != nil {
		return err
	}

	err = accInput.CPF.ValidateInput()
	if err != nil {
		return err
	}

	return nil
}

// validateInputEmpty validates if the user has filled the required fields
func (accInput *AccountInput) validateInputEmpty() error {
	if accInput.Name == "" || accInput.CPF == "" || accInput.Secret == "" {
		return entities.ErrEmptyInput
	}
	return nil
}

// removeBlankSpaces removes blank spaces of account fields
func (accInput *AccountInput) removeBlankSpaces() {
	accInput.Name = strings.TrimSpace(accInput.Name)
	accInput.CPF = vos.CPF(strings.TrimSpace(accInput.CPF.String()))
	accInput.Secret = strings.TrimSpace(accInput.Secret)
}
