package usecase

import (
	"context"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"strings"
	"time"
)

const balanceInit vos.Currency = 10000

//AccountInput represents information necessary to create a bank account
type AccountInput struct {
	Name   string  `json:"name"`
	CPF    vos.CPF `json:"cpf"`
	Secret string  `json:"secret"`
}

//CreateAccount validates and handles user input and creates a formatted account,
//then calls the function to insert the account into the database
func (accUseCase *accountUseCase) CreateAccount(ctx context.Context, accInput *AccountInput) (*models.AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	accID := vos.NewUUID()
	account := &models.Account{
		ID:        accID,
		Name:      accInput.Name,
		CPF:       accInput.CPF,
		Secret:    accInput.Secret,
		Balance:   balanceInit,
		CreatedAt: time.Now().Truncate(time.Second),
	}

	err := account.GetHashSecret()
	if err != nil {
		accUseCase.log.WithError(err).Println(domainerr.ErrUnexpected)
		return nil, domainerr.ErrUnexpected
	}

	account.Balance.ConvertToCents()
	err = accUseCase.accountRepo.CreateAccount(ctx, account)

	if err != nil {
		return nil, err
	}
	return account.GetAccOutput(), nil
}

//ValidateAccountInput validates account input and returns if occurred an error
func (accInput *AccountInput) ValidateAccountInput() error {
	accInput.removeBlankSpaces()

	err := accInput.validateInputEmpty()
	if err != nil {
		return err
	}

	err = accInput.validateSecretLen()
	if err != nil {
		return err
	}

	err = accInput.CPF.ValidateInput()
	if err != nil {
		return err
	}

	return nil
}

//validateInputEmpty validates if the user has filled the required fields
func (accInput *AccountInput) validateInputEmpty() error {
	if accInput.Name == "" || accInput.CPF == "" || accInput.Secret == "" {
		return domainerr.ErrEmptyInput
	}
	return nil
}

//secretLen validates the secret length
func (accInput *AccountInput) validateSecretLen() error {
	if len(accInput.Secret) < 8 {
		return domainerr.ErrSmallSecret
	}
	return nil
}

//removesBlankSpaces removes blank spaces of account fields
func (accInput *AccountInput) removeBlankSpaces() {
	accInput.Name = strings.TrimSpace(accInput.Name)
	accInput.CPF = vos.CPF(strings.TrimSpace(accInput.CPF.String()))
	accInput.Secret = strings.TrimSpace(accInput.Secret)
}
