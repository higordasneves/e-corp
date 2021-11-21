package controller

import (
	"encoding/json"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/responses"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type AccountController interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
}

type accountController struct {
	accUseCase usecase.AccountUseCase
	log        *logrus.Logger
}

func NewAccountController(accUseCase usecase.AccountUseCase, log *logrus.Logger) AccountController {
	return &accountController{accUseCase: accUseCase, log: log}
}

// CreateAccount CreateAccount reads HTTP POST request to create an account and returns a response
func (accController accountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err, accController.log)
		return
	}

	var accountInput usecase.AccountInput
	if err = json.Unmarshal(bodyRequest, &accountInput); err != nil {
		responses.Error(w, http.StatusBadRequest, err, accController.log)
		return
	}

	err = accountInput.ValidateAccountInput()
	if err != nil {
		accountInput.Secret = "######"
		responses.Error(w, http.StatusBadRequest, err, accController.log)
		return
	}

	account, err := accController.accUseCase.CreateAccount(r.Context(), accountInput)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err, accController.log)
		return
	}
	accOutput := account.GetAccOutput()
	responses.JSON(w, http.StatusCreated, accOutput, accController.log)
}
