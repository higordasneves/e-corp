package controller

import (
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"

	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
)

// GetBalance returns the current balance of the account.
// It returns NotFound error if the account not exists.
func (accController AccountController) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(mux.Vars(r)["account_id"])
	balance, err := accController.accUseCase.GetBalance(r.Context(), id)
	if err != nil {
		reponses.HandleError(w, err, accController.log)
		return
	}

	balanceResponse := map[string]int{"balance": balance}
	reponses.SendResponse(w, http.StatusOK, balanceResponse, accController.log)
}
