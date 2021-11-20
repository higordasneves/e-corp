package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
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
		Balance:   accInput.Balance,
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

func ValidateAccountInput(accInput *AccountInput) error {
	err := accInput.InputEmpty()
	if err != nil {
		return err
	}

	err = accInput.SecretLen()
	if err != nil {
		return err
	}

	err = accInput.CPFLen()
	if err != nil {
		return err
	}
	return nil
}

func (accInput *AccountInput) InputEmpty() error {
	if accInput.Name == "" || accInput.CPF == "" || accInput.Secret == "" {
		return ErrEmptyInput
	}
	return nil
}

func (accInput *AccountInput) SecretLen() error {
	if len(accInput.Secret) < 8 {
		return ErrSmallSecret
	}
	return nil
}

func (accInput *AccountInput) CPFLen() error {
	if len(accInput.CPF) != 8 {
		return ErrCPFLen
	}
	return nil
}

func (accInput *AccountInput) CPFFormat() error {
	for _, v := range accInput.CPF {
		if !unicode.IsDigit(v) {
			return ErrCPFFormat
		}
	}
	return nil
}
