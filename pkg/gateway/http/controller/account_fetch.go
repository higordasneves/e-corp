package controller

import (
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"net/http"
)

// FetchAccounts reads HTTP GET request for accounts and sends response with account list or error
func (accController accountController) FetchAccounts(w http.ResponseWriter, r *http.Request) {
	accList, err := accController.accUseCase.FetchAccounts(r.Context())
	if err != nil {
		interpreter.HandleError(w, err, accController.log)
		return
	}
	interpreter.SendResponse(w, http.StatusOK, accList, accController.log)
}
