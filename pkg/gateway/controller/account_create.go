package controller

import (
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
	"net/http"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

// CreateAccount reads HTTP POST request to create an account and returns a response
func (accController AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var accountInput usecase.CreateAccountInput
	if err := requests.ReadRequestBody(r, &accountInput); err != nil {
		reponses.HandleError(w, err, accController.log)
		return
	}

	account, err := accController.accUseCase.CreateAccount(r.Context(), &accountInput)

	if err != nil {
		reponses.HandleError(w, err, accController.log)
		return
	}
	reponses.SendResponse(w, http.StatusCreated, account, accController.log)
}
