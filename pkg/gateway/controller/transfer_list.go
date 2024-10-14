package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

type ListTransfersResponse struct {
	Transfers []ListTransfersResponseItem `json:"transfers"`
}

// ListTransfersResponseItem represents a banking transfer.
type ListTransfersResponseItem struct {
	ID                   uuid.UUID `json:"id"`
	AccountOriginID      uuid.UUID `json:"account_origin_id"`
	AccountDestinationID uuid.UUID `json:"account_destination_id"`
	Amount               int       `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

// ListTransfers lists all the transfers sent or received by the account in desc order.
// Returns not found error if the account not exists.
// @Summary List Transfers
// @Description Lists all the transfers sent or received by the account in desc order.
// @Description It returns not found error if the account not exists.
// @Description The account id is obtained from the subject.
// @Tags Transfers
// @Accept json
// @Produce json
// @Success 200 {object} ListTransfersResponse "Transfers list"
// @Failure 404 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/transfers [get]
func (tController TransferController) ListTransfers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountOriginID := fmt.Sprint(r.Context().Value("subject"))
	id, err := uuid.FromString(accountOriginID)
	if err != nil {
		HandleError(ctx, w, fmt.Errorf("unexpected error when parsing the account id: %w", err))
		return
	}

	ucOutput, err := tController.tUseCase.ListAccountTransfers(r.Context(), usecase.ListAccountTransfersInput{
		AccountID: id,
	})
	if err != nil {
		HandleError(ctx, w, err)
		return
	}

	resp := make([]ListTransfersResponseItem, 0, len(ucOutput.Transfers))
	for _, transfer := range ucOutput.Transfers {
		resp = append(resp, ListTransfersResponseItem(transfer))
	}

	SendResponse(ctx, w, http.StatusOK, ListTransfersResponse{Transfers: resp})
}
