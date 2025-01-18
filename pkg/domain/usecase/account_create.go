package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

//go:generate moq -stub -pkg mocks -out mocks/account_create.go . CreateAccountUCBroker

type CreateAccountUCRepository interface {
	CreateAccount(ctx context.Context, acc entities.Account) error
}

type CreateAccountUCBroker interface {
	NotifyAccountCreation(ctx context.Context, account entities.Account) error
}

type CreateAccountUC struct {
	R CreateAccountUCRepository
	B CreateAccountUCBroker
}

func NewCreateAccountUC(accountRepo CreateAccountUCRepository, broker CreateAccountUCBroker) CreateAccountUC {
	return CreateAccountUC{R: accountRepo, B: broker}
}

// CreateAccountInput represents information necessary to create a bank account.
type CreateAccountInput struct {
	Name     string
	Document string
	Secret   string
}

type CreateAccountOutput struct {
	Account entities.Account
}

// CreateAccount validates the input and creates an account.
// Returns domain.ErrInvalidParameter if:
// - the account name is not filled;
// - the number of characters of the document is not valid;
// - the format of the document is not valid;
// - the number of the characters of the secret is less than the minimum;
// - the account already exists.
func (accUseCase CreateAccountUC) CreateAccount(ctx context.Context, input CreateAccountInput) (CreateAccountOutput, error) {
	input = input.removeBlankSpaces()
	if input.Name == "" {
		return CreateAccountOutput{}, fmt.Errorf("%w (name): required field", domain.ErrInvalidParameter)
	}

	document, err := vos.NewDocument(input.Document)
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("%w (document): %w", domain.ErrInvalidParameter, err)
	}

	secret, err := vos.NewSecret(input.Secret)
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("%w (secret): %w", domain.ErrInvalidParameter, err)
	}

	account := entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      input.Name,
		Document:  document,
		Secret:    secret,
		Balance:   0,
		CreatedAt: time.Now().Truncate(time.Second),
	}

	err = accUseCase.R.CreateAccount(ctx, account)
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("creating account in the database: %w", err)
	}

	err = accUseCase.B.NotifyAccountCreation(ctx, account)
	if err != nil {
		return CreateAccountOutput{}, fmt.Errorf("notifying account creation in the broker: %w", err)
	}

	return CreateAccountOutput{account}, nil
}

// removeBlankSpaces removes blank spaces of account fields
func (accInput CreateAccountInput) removeBlankSpaces() CreateAccountInput {
	accInput.Name = strings.TrimSpace(accInput.Name)
	accInput.Document = strings.TrimSpace(accInput.Document)
	accInput.Secret = strings.TrimSpace(accInput.Secret)

	return accInput
}
