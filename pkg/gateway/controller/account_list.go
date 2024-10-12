package controller

import (
	"github.com/gofrs/uuid/v5"
	"net/http"
	"time"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
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
	// todo: add cursor.
	ucOutput, err := accController.accUseCase.ListAccounts(r.Context(), usecase.ListAccountsInput{
		PageSize: 100,
	})
	if err != nil {
		reponses.HandleError(w, err, accController.log)
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

	response := ListAccountsResponse{
		Accounts: responseItems,
		// todo: add cursor.
		NextPage: "",
	}

	reponses.SendResponse(w, http.StatusOK, response, accController.log)
}
