package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"strings"
	"time"
)

//TransferInput represents information necessary to transfer money between bank accounts
type TransferInput struct {
	AccountOriginID      string
	AccountDestinationID string `json:"destinationID"`
	Amount               int    `json:"amount"`
}

//Transfer creates a transfer and updates account balances of the clients
func (tUseCase transferUseCase) Transfer(ctx context.Context, transferInput *TransferInput) (*entities.Transfer, error) {
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

	ctxChan := make(chan context.Context)
	errChan := make(chan error)

	go func() {
		ctxWithValue, ok := <-ctxChan
		if !ok {
			return
		}
		errChan <- tUseCase.transferRepo.CreateTransfer(ctxWithValue, transfer)
	}()

	err = tUseCase.transferRepo.Transfer(ctx, ctxChan, errChan)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

//ValidateInput validates transfer input
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

//validateBalance validates if transfer balance is greater than zero
func (transferInput *TransferInput) validateBalance() error {
	if transferInput.Amount <= 0 {
		return entities.ErrTransferAmount
	}
	return nil
}

//removesBlankSpaces removes blank spaces of transfer fields
func (transferInput *TransferInput) removeBlankSpaces() {
	transferInput.AccountOriginID = strings.TrimSpace(transferInput.AccountOriginID)
	transferInput.AccountDestinationID = strings.TrimSpace(transferInput.AccountDestinationID)
}
