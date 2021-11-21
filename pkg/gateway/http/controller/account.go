package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/responses"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type AccountController interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
	FetchAccounts(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
}

type accountController struct {
	accUseCase usecase.AccountUseCase
	log        *logrus.Logger
}

func NewAccountController(accUseCase usecase.AccountUseCase, log *logrus.Logger) AccountController {
	return &accountController{accUseCase: accUseCase, log: log}
}

// CreateAccount reads HTTP POST request to create an account and returns a response
func (accController accountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendResponse(w, http.StatusCreated, responses.ErrorJSON(err), accController.log)
		return
	}

	var accountInput usecase.AccountInput
	if err = json.Unmarshal(bodyRequest, &accountInput); err != nil {
		responses.SendResponse(w, http.StatusCreated, responses.ErrorJSON(err), accController.log)
		return
	}

	err = accountInput.ValidateAccountInput()
	if err != nil {
		accountInput.Secret = "######"
		responses.SendResponse(w, http.StatusCreated, responses.ErrorJSON(err), accController.log)
		return
	}

	account, err := accController.accUseCase.CreateAccount(r.Context(), accountInput)
	if err != nil {
		responses.SendResponse(w, http.StatusCreated, responses.ErrorJSON(err), accController.log)
		return
	}
	accOutput := account.GetAccOutput()
	responses.SendResponse(w, http.StatusCreated, accOutput, accController.log)
}

// FetchAccounts reads HTTP GET request for accounts and sends response with account list or error
func (accController accountController) FetchAccounts(w http.ResponseWriter, r *http.Request) {
	accList, err := accController.accUseCase.FetchAccounts(r.Context())
	if err != nil {
		responses.SendResponse(w, http.StatusInternalServerError, responses.ErrorJSON(err), accController.log)
		return
	}

	responses.SendResponse(w, http.StatusCreated, accList, accController.log)

}

func (accController accountController) GetBalance(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["account_id"]

	balance, err := accController.accUseCase.GetBalance(r.Context(), vos.AccountID(id))

	if err == postgres.ErrAccNotFound {
		responses.SendResponse(w, http.StatusBadRequest, responses.ErrorJSON(err), accController.log)
		return
	}

	if err != nil {
		responses.SendResponse(w, http.StatusInternalServerError, responses.ErrorJSON(err), accController.log)
		return
	}

	balanceResponse := map[string]*vos.Currency{"balance": balance}
	responses.SendResponse(w, http.StatusCreated, balanceResponse, accController.log)

}
