package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
)

type TransferRequest struct {
	AccountDestinationID uuid.UUID `json:"destination_id"`
	Amount               int       `json:"amount"`
}

type TransferResponse struct {
	ID                   uuid.UUID `json:"id"`
	AccountOriginID      uuid.UUID `json:"account_origin_id"`
	AccountDestinationID uuid.UUID `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

// Transfer creates a transfer and updates the balance of the destination and origin accounts.
// Returns bad request error if:
// - The AccountOriginID is equal to AccountDestinationID.
// - The amount is less than or equal to zero.
// - The origin accounts doesn't have enough funds to complete the transfer.
// Returns not found error if the destination account not exists.
func (tController TransferController) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TransferRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	accountOriginID := uuid.FromStringOrNil(fmt.Sprint(r.Context().Value("subject")))

	ucOutput, err := tController.tUseCase.Transfer(r.Context(), usecase.TransferInput{
		AccountOriginID:      accountOriginID,
		AccountDestinationID: req.AccountDestinationID,
		Amount:               req.Amount,
	})
	if err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	reponses.SendResponse(ctx, w, http.StatusCreated, TransferResponse(ucOutput.Transfer))
}
