package controller

import (
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
)

type GetBalanceResponse struct {
	// Balance represents the balance of the account.
	Balance int `json:"balance"`
}

// GetBalance returns the current balance of the account.
// It returns NotFound error if the account not exists.
// @Summary Get Balance
// @Description Returns the current balance of the account.
// @Description It returns NotFound error if the account not exists.
// @Tags Accounts
// @Param account_id path string true "Account ID"
// @Accept json
// @Produce json
// @Success 200 {object} GetBalanceResponse "Account Balance"
// @Failure 404 {object} ErrorResponse "Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/accounts/{account_id}/balance [GET]
func (accController AccountController) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := uuid.FromStringOrNil(mux.Vars(r)["account_id"])
	balance, err := accController.accUseCase.GetBalance(r.Context(), id)
	if err != nil {
		HandleError(ctx, w, err)
		return
	}

	balanceResponse := GetBalanceResponse{Balance: balance}
	SendResponse(ctx, w, http.StatusOK, balanceResponse)
}
