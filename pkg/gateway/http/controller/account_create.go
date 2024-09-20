package controller

import (
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"net/http"
)

// CreateAccount reads HTTP POST request to create an account and returns a response
func (accController AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var accountInput usecase.AccountInput
	if err := interpreter.ReadRequestBody(r, &accountInput); err != nil {
		interpreter.HandleError(w, err, accController.log)
		return
	}

	account, err := accController.accUseCase.CreateAccount(r.Context(), &accountInput)

	if err != nil {
		interpreter.HandleError(w, err, accController.log)
		return
	}
	interpreter.SendResponse(w, http.StatusCreated, account, accController.log)
}
