package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

// TransferInput represents information necessary to transfer money between bank accounts
type TransferInput struct {
	AccountOriginID      string
	AccountDestinationID string `json:"destinationID"`
	Amount               int    `json:"amount"`
}

// Transfer creates a transfer and updates account balances of the clients
func (tUseCase TransferUseCase) Transfer(ctx context.Context, transferInput *TransferInput) (*entities.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	transferInput.removeBlankSpaces()
	err := transferInput.ValidateInput()
	if err != nil {
		return nil, err
	}

	transfer := &entities.Transfer{
		ID:                   vos.NewUUID(),
		AccountOriginID:      vos.UUID(transferInput.AccountOriginID),
		AccountDestinationID: vos.UUID(transferInput.AccountDestinationID),
		Amount:               transferInput.Amount,
		CreatedAt:            time.Now().Truncate(time.Second),
	}

	err = tUseCase.validateAccounts(ctx, transfer)
	if err != nil {
		return nil, err
	}

	ctxChan := make(chan context.Context)
	errChan := make(chan error)

	go func() {
		ctxWithValue, ok := <-ctxChan
		if !ok {
			return
		}

		err = tUseCase.repo.CreateTransfer(ctxWithValue, transfer)
		if err != nil {
			errChan <- err
			return
		}

		err = tUseCase.repo.UpdateBalance(ctxWithValue, transfer.AccountOriginID, -transfer.Amount)
		if err != nil {
			errChan <- err
			return
		}

		err = tUseCase.repo.UpdateBalance(ctxWithValue, transfer.AccountDestinationID, transfer.Amount)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	err = tUseCase.repo.PerformTransaction(ctx, ctxChan, errChan)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

// validateAccounts validates existence of the accounts involved and balance sufficiency
func (tUseCase TransferUseCase) validateAccounts(ctx context.Context, transfer *entities.Transfer) error {
	_, err := tUseCase.repo.GetBalance(ctx, transfer.AccountDestinationID)
	if err != nil {
		if errors.Is(err, entities.ErrAccNotFound) {
			return fmt.Errorf("destination %w", err)
		}
		return err
	}

	originBalance, err := tUseCase.repo.GetBalance(ctx, transfer.AccountOriginID)
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
	err := vos.IsValidUUID(transferInput.AccountOriginID)
	if err != nil {
		return entities.ErrOriginAccID
	}

	err = vos.IsValidUUID(transferInput.AccountDestinationID)
	if err != nil {
		return entities.ErrDestAccID
	}

	if transferInput.AccountOriginID == transferInput.AccountDestinationID {
		return entities.ErrSelfTransfer
	}

	err = transferInput.validateBalance()
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

// removesBlankSpaces removes blank spaces of transfer fields
func (transferInput *TransferInput) removeBlankSpaces() {
	transferInput.AccountOriginID = strings.TrimSpace(transferInput.AccountOriginID)
	transferInput.AccountDestinationID = strings.TrimSpace(transferInput.AccountDestinationID)
}
