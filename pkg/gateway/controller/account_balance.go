package controller

import (
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"net/http"
)

func (accController AccountController) GetBalance(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := uuid.FromStringOrNil(params["account_id"])

	balance, err := accController.accUseCase.GetBalance(r.Context(), id)
	if err != nil {
		reponses.HandleError(w, err, accController.log)
		return
	}

	balanceResponse := map[string]int{"balance": balance}
	reponses.SendResponse(w, http.StatusOK, balanceResponse, accController.log)
}
