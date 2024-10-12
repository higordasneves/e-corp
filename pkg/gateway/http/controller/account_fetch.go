package controller

import (
	"net/http"

	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/reponses"
)

// FetchAccounts reads HTTP GET request for accounts and sends response with account list or error
func (accController AccountController) FetchAccounts(w http.ResponseWriter, r *http.Request) {
	accList, err := accController.accUseCase.FetchAccounts(r.Context())
	if err != nil {
		reponses.HandleError(w, err, accController.log)
		return
	}
	reponses.SendResponse(w, http.StatusOK, accList, accController.log)
}
