package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type TransferUCRepository interface {
	GetBalance(ctx context.Context, id uuid.UUID) (int, error)
	GetAccountByDocument(ctx context.Context, cpf vos.Document) (entities.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, transactionAmount int) error

	CreateTransfer(ctx context.Context, transfer entities.Transfer) error

	BeginTX(ctx context.Context) (context.Context, error)
	CommitTX(ctx context.Context) error
	RollbackTX(ctx context.Context) error
}

type TransferUC struct {
	R TransferUCRepository
}

func NewTransferUC(r TransferUCRepository) TransferUC {
	return TransferUC{R: r}
}

// TransferInput represents information necessary to transfer money between bank accounts
type TransferInput struct {
	AccountOriginID      uuid.UUID
	AccountDestinationID uuid.UUID `json:"destinationID"`
	Amount               int       `json:"amount"`
}

type TransferOutput struct {
	Transfer entities.Transfer
}

// Transfer creates a transfer and updates the balance of the destination and origin accounts.
// Returns domain.ErrInvalidParameter if:
// - The AccountOriginID is equal to AccountDestinationID.
// - The amount is less than or equal to zero.
// - The origin accounts doesn't have enough funds to complete the transfer.
// Returns domain.ErrNotFound if the origin or destination account not exists.
func (tUseCase TransferUC) Transfer(ctx context.Context, input TransferInput) (TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	err := ValidateTransferInput(input)
	if err != nil {
		return TransferOutput{}, fmt.Errorf("%w: %w", domain.ErrInvalidParameter, err)
	}

	transfer := entities.Transfer{
		ID:                   uuid.Must(uuid.NewV7()),
		AccountOriginID:      input.AccountOriginID,
		AccountDestinationID: input.AccountDestinationID,
		Amount:               input.Amount,
		CreatedAt:            time.Now().Truncate(time.Second),
	}

	err = tUseCase.validate(ctx, transfer)
	if err != nil {
		return TransferOutput{}, fmt.Errorf("%w: %w", domain.ErrInvalidParameter, err)
	}

	ctx, err = tUseCase.R.BeginTX(ctx)
	if err != nil {
		return TransferOutput{}, fmt.Errorf("starting transaction: %w", err)
	}
	defer tUseCase.R.RollbackTX(ctx) // nolint:errcheck

	err = tUseCase.R.CreateTransfer(ctx, transfer)
	if err != nil {
		return TransferOutput{}, fmt.Errorf("error creating transfer: %w", err)
	}

	err = tUseCase.R.UpdateBalance(ctx, transfer.AccountOriginID, -transfer.Amount)
	if err != nil {
		return TransferOutput{}, fmt.Errorf("error updating origin account balance: %w", err)
	}

	err = tUseCase.R.UpdateBalance(ctx, transfer.AccountDestinationID, transfer.Amount)
	if err != nil {
		return TransferOutput{}, fmt.Errorf("error updating destination account balance: %w", err)
	}

	if err = tUseCase.R.CommitTX(ctx); err != nil {
		return TransferOutput{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return TransferOutput{transfer}, nil
}

// ValidateTransferInput validates transfer input.
// Returns domain.ErrInvalidParameter if the AccountOriginID is equal to AccountDestinationID.
// Returns domain.ErrInvalidParameter if the amount is less than or equal to zero.
func ValidateTransferInput(i TransferInput) error {
	if i.AccountOriginID == i.AccountDestinationID {
		return fmt.Errorf("%w: the destination account must be different from the origin account", domain.ErrInvalidParameter)
	}

	if i.Amount <= 0 {
		return fmt.Errorf("%w: invalid transfer amount, the amount must be greater than 0", domain.ErrInvalidParameter)
	}

	return nil
}

// validate validates existence of the accounts involved and balance sufficiency.
// Returns domain.ErrNotFound if the origin or destination account not exists.
// Returns domain.ErrInvalidParameter if the origin accounts doesn't have enough funds to complete the transfer.
func (tUseCase TransferUC) validate(ctx context.Context, transfer entities.Transfer) error {
	originBalance, err := tUseCase.R.GetBalance(ctx, transfer.AccountOriginID)
	if err != nil {
		return fmt.Errorf("getting origin account balance: %w", err)
	}

	// just checking if the destination account exists.
	_, err = tUseCase.R.GetBalance(ctx, transfer.AccountDestinationID)
	if err != nil {
		return fmt.Errorf("getting destination account balance: %w", err)
	}

	if transfer.Amount > originBalance {
		return fmt.Errorf("%w: insufficient funds", domain.ErrInvalidParameter)
	}

	return nil
}
