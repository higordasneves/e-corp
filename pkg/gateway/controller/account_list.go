package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"github.com/higordasneves/e-corp/utils/pagination"
)

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

	var ucInput usecase.ListAccountsInput
	if t := r.URL.Query().Get("page_token"); t != "" {
		err := pagination.Extract(t, &ucInput)
		if err != nil {
			reponses.HandleError(ctx, w, fmt.Errorf("%w: invalid page token", domain.ErrInvalidParameter))
		}
	} else {
		pageSize := r.URL.Query().Get("page_size")
		if pageSize == "" {
			pageSize = "0"
		}
		i, err := strconv.Atoi(pageSize)
		if err != nil {
			reponses.HandleError(ctx, w, fmt.Errorf("%w: converting page size to int", domain.ErrInvalidParameter))
			return
		}

		idsString := strings.Split(r.URL.Query().Get("ids"), ",")
		accountIDs := make([]uuid.UUID, 0, len(idsString))
		for _, id := range idsString {
			accountID, err := uuid.FromString(id)
			if err != nil {
				reponses.HandleError(ctx, w, fmt.Errorf("%w: invalid account id", domain.ErrInvalidParameter))
				return
			}
			accountIDs = append(accountIDs, accountID)
		}

		ucInput.PageSize = pagination.ValidatePageSize(uint32(i)) // nolint:gosec
		ucInput.IDs = accountIDs
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
