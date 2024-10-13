package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
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
func (tController TransferController) ListTransfers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accountOriginID := fmt.Sprint(r.Context().Value("subject"))
	id, err := uuid.FromString(accountOriginID)
	if err != nil {
		reponses.HandleError(ctx, w, fmt.Errorf("unexpected error when parsing the account id: %w", err))
		return
	}

	ucOutput, err := tController.tUseCase.ListAccountTransfers(r.Context(), usecase.ListAccountTransfersInput{
		AccountID: id,
	})
	if err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	resp := make([]ListTransfersResponseItem, 0, len(ucOutput.Transfers))
	for _, transfer := range ucOutput.Transfers {
		resp = append(resp, ListTransfersResponseItem(transfer))
	}

	reponses.SendResponse(ctx, w, http.StatusOK, ListTransfersResponse{Transfers: resp})
}
