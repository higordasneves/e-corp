package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

// TransferInput represents information necessary to transfer money between bank accounts
type TransferInput struct {
	AccountOriginID      uuid.UUID
	AccountDestinationID uuid.UUID `json:"destinationID"`
	Amount               int       `json:"amount"`
}

// Transfer creates a transfer and updates account balances of the clients
func (tUseCase TransferUseCase) Transfer(ctx context.Context, transferInput *TransferInput) (*entities.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	err := transferInput.ValidateInput()
	if err != nil {
		return nil, err
	}

	transfer := entities.Transfer{
		ID:                   uuid.Must(uuid.NewV7()),
		AccountOriginID:      transferInput.AccountOriginID,
		AccountDestinationID: transferInput.AccountDestinationID,
		Amount:               transferInput.Amount,
		CreatedAt:            time.Now().Truncate(time.Second),
	}

	err = tUseCase.validateAccounts(ctx, transfer)
	if err != nil {
		return nil, err
	}

	ctx, err = tUseCase.R.BeginTX(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer tUseCase.R.RollbackTX(ctx)

	err = tUseCase.R.CreateTransfer(ctx, transfer)
	if err != nil {
		return nil, fmt.Errorf("error creating transfer: %w", err)
	}

	err = tUseCase.R.UpdateBalance(ctx, transfer.AccountOriginID, -transfer.Amount)
	if err != nil {
		return nil, fmt.Errorf("error updating origin account balance: %w", err)
	}

	err = tUseCase.R.UpdateBalance(ctx, transfer.AccountDestinationID, transfer.Amount)
	if err != nil {
		return nil, fmt.Errorf("error updating destination account balance: %w", err)
	}

	if err = tUseCase.R.CommitTX(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &transfer, nil
}

// validateAccounts validates existence of the accounts involved and balance sufficiency
func (tUseCase TransferUseCase) validateAccounts(ctx context.Context, transfer entities.Transfer) error {
	_, err := tUseCase.R.GetBalance(ctx, transfer.AccountDestinationID)
	if err != nil {
		if errors.Is(err, entities.ErrAccNotFound) {
			return fmt.Errorf("destination %w", err)
		}
		return err
	}

	originBalance, err := tUseCase.R.GetBalance(ctx, transfer.AccountOriginID)
	if err != nil {
		if errors.Is(err, entities.ErrAccNotFound) {
			return fmt.Errorf("origin %w", err)
		}
		return err
	}

	if transfer.Amount > originBalance {
		return entities.ErrTransferInsufficientFunds
	}
	return nil
}

// ValidateInput validates transfer input
func (transferInput *TransferInput) ValidateInput() error {
	if transferInput.AccountOriginID == transferInput.AccountDestinationID {
		return entities.ErrSelfTransfer
	}

	err := transferInput.validateBalance()
	if err != nil {
		return err
	}

	return nil
}

// validateBalance validates if transfer balance is greater than zero
func (transferInput *TransferInput) validateBalance() error {
	if transferInput.Amount <= 0 {
		return entities.ErrTransferAmount
	}
	return nil
}
