package usecase

import (
	"context"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"strings"
	"time"
	"unicode"
)

const balanceInit vos.Currency = 10000

//AccountInput represents information necessary to create a bank account
type AccountInput struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
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

	err := accInput.inputEmpty()
	if err != nil {
		return err
	}

	err = accInput.secretLen()
	if err != nil {
		return err
	}

	err = accInput.cpfFormat()
	if err != nil {
		return err
	}

	err = accInput.cpfLen()
	if err != nil {
		return err
	}

	return nil
}

//inputEmpty validates if the user has filled the required fields
func (accInput *AccountInput) inputEmpty() error {
	if accInput.Name == "" || accInput.CPF == "" || accInput.Secret == "" {
		return domainerr.ErrEmptyInput
	}
	return nil
}

//secretLen validates the secret length
func (accInput *AccountInput) secretLen() error {
	if len(accInput.Secret) < 8 {
		return domainerr.ErrSmallSecret
	}
	return nil
}

//cpfLen validates the CPF length
func (accInput *AccountInput) cpfLen() error {
	if len(accInput.CPF) != 11 {
		return domainerr.ErrCPFLen
	}
	return nil
}

//cpfLen validates if the CPF has only numbers
func (accInput *AccountInput) cpfFormat() error {
	for _, v := range accInput.CPF {
		if !unicode.IsDigit(v) {
			return domainerr.ErrCPFFormat
		}
	}
	return nil
}

//removesBlankSpaces removes blank spaces of account fields
func (accInput *AccountInput) removeBlankSpaces() {
	accInput.Name = strings.TrimSpace(accInput.Name)
	accInput.CPF = strings.TrimSpace(accInput.CPF)
	accInput.Secret = strings.TrimSpace(accInput.Secret)
}
