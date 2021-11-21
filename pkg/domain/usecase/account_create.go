package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"strings"
	"time"
	"unicode"
)

var (
	ErrEmptyInput  = errors.New("the name, document and password fields are required")
	ErrSmallSecret = errors.New("the password must be at least 8 characters long")
	ErrCPFLen      = errors.New("the CPF must be 11 characters long")
	ErrCPFFormat   = errors.New("the CPF must contain only numbers")
)

type AccountInput struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"secret"`
	Balance float64 `json:"balance"`
}

//CreateAccount validates and handles user input and creates a formatted account,
//then calls the function to insert the account into the database
func (accUseCase *accountUseCase) CreateAccount(ctx context.Context, accInput AccountInput) (*models.Account, error) {
	accID := newAccID()

	account := &models.Account{ID: accID,
		Name:      accInput.Name,
		CPF:       accInput.CPF,
		Secret:    accInput.Secret,
		Balance:   vos.Currency(accInput.Balance),
		CreatedAt: time.Now().Truncate(time.Second),
	}

	account.GetHashSecret()

	err := accUseCase.accountRepo.CreateAccount(ctx, account)

	if err != nil {
		return nil, err
	}
	return account, nil
}

// newAccID gets uuid using google lib
func newAccID() vos.AccountID {
	return vos.AccountID(uuid.NewString())
}

//ValidateAccountInput validates account input and returns if occurred an error
func (accInput *AccountInput) ValidateAccountInput() error {
	accInput.removesBlankSpaces()

	err := accInput.inputEmpty()
	if err != nil {
		return err
	}

	err = accInput.secretLen()
	if err != nil {
		return err
	}

	err = accInput.cpfLen()
	if err != nil {
		return err
	}

	err = accInput.cpfFormat()
	if err != nil {
		return err
	}

	return nil
}

//inputEmpty validates if the user has filled the required fields
func (accInput *AccountInput) inputEmpty() error {
	if accInput.Name == "" || accInput.CPF == "" || accInput.Secret == "" {
		return ErrEmptyInput
	}
	return nil
}

//secretLen validates the secret length
func (accInput *AccountInput) secretLen() error {
	if len(accInput.Secret) < 8 {
		return ErrSmallSecret
	}
	return nil
}

//cpfLen validates the CPF length
func (accInput *AccountInput) cpfLen() error {
	if len(accInput.CPF) != 11 {
		return ErrCPFLen
	}
	return nil
}

//cpfLen validates if the CPF has only numbers
func (accInput *AccountInput) cpfFormat() error {
	for _, v := range accInput.CPF {
		if !unicode.IsDigit(v) {
			return ErrCPFFormat
		}
	}
	return nil
}

//removesBlankSpaces removes blank spaces of account fields
func (accInput *AccountInput) removesBlankSpaces() {
	accInput.Name = strings.TrimSpace(accInput.Name)
	accInput.CPF = strings.TrimSpace(accInput.CPF)
	accInput.Secret = strings.TrimSpace(accInput.Secret)
}
