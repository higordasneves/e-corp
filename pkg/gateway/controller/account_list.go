package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
	"github.com/higordasneves/e-corp/utils/pagination"
)

type ListAccountsRequest struct {
	IDs       []uuid.UUID `json:"ids"`
	PageSize  int         `json:"page_size"`
	PageToken string      `json:"page_token"`
}

type ListAccountsResponse struct {
	Accounts []ListAccountsResponseItem `json:"accounts"`
	NextPage string                     `json:"next_page"`
}

type ListAccountsResponseItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// ListAccounts Lists accounts by filtering the IDs provided in the input.
func (accController AccountController) ListAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ListAccountsRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	var ucInput usecase.ListAccountsInput
	if req.PageToken != "" {
		err := pagination.Extract(req.PageToken, &ucInput)
		if err != nil {
			reponses.HandleError(ctx, w, fmt.Errorf("%w: invalid page token", domain.ErrInvalidParameter))
		}
	} else {
		ucInput.PageSize = pagination.ValidatePageSize(uint32(req.PageSize))
		ucInput.IDs = req.IDs
	}

	ucOutput, err := accController.accUseCase.ListAccounts(r.Context(), ucInput)
	if err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	responseItems := make([]ListAccountsResponseItem, 0, len(ucOutput.Accounts))
	for _, acc := range ucOutput.Accounts {
		responseItems = append(responseItems, ListAccountsResponseItem{
			ID:        acc.ID,
			Name:      acc.Name,
			Document:  acc.Document.String(),
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
		})
	}

	var nextPageToken string
	if ucOutput.NextPage != nil {
		v, err := pagination.NewToken(*ucOutput.NextPage)
		if err != nil {
			reponses.HandleError(ctx, w, errors.New("unexpect error"))
		}
		nextPageToken = v
	}

	response := ListAccountsResponse{
		Accounts: responseItems,
		NextPage: nextPageToken,
	}

	reponses.SendResponse(ctx, w, http.StatusOK, response)
}
