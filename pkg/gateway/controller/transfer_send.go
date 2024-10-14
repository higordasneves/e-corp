package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
)

type TransferRequest struct {
	AccountDestinationID uuid.UUID `json:"destination_id"`
	// Amount is the amount of the transfer. It must be positive.
	Amount int `json:"amount"`
}

type TransferResponse struct {
	ID                   uuid.UUID `json:"id"`
	AccountOriginID      uuid.UUID `json:"account_origin_id"`
	AccountDestinationID uuid.UUID `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

// Transfer creates a transfer and updates the balance of the destination and origin accounts.
// @Summary Send Transfer
// @Description Creates a transfer and updates the balance of the destination and origin accounts.
// @Description The origin account id is obtained from the subject.
// @Description It returns not found error if the destination account not exists.
// @Description It returns bad request error if:
// @Description - The AccountOriginID is equal to AccountDestinationID.
// @Description - The amount is less than or equal to zero.
// @Description - The origin accounts doesn't have enough funds to complete the transfer.
// @Tags Transfers
// @Param Body body TransferRequest true "Request body"
// @Accept json
// @Produce json
// @Success 200 {object} TransferResponse "Transfer Created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/transfers [post]
func (tController TransferController) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TransferRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		HandleError(ctx, w, err)
		return
	}

	accountOriginID := uuid.FromStringOrNil(fmt.Sprint(r.Context().Value("subject")))

	ucOutput, err := tController.tUseCase.Transfer(r.Context(), usecase.TransferInput{
		AccountOriginID:      accountOriginID,
		AccountDestinationID: req.AccountDestinationID,
		Amount:               req.Amount,
	})
	if err != nil {
		HandleError(ctx, w, err)
		return
	}

	SendResponse(ctx, w, http.StatusCreated, TransferResponse(ucOutput.Transfer))
}
