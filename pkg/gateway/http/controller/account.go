package controller

import (
	"encoding/json"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
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
		accController.log.WithError(err).Warn()
	}

	var accountInput usecase.AccountInput
	if err = json.Unmarshal(bodyRequest, &accountInput); err != nil {
		accController.log.WithError(err).Warn()
	}

	account, err := accController.accUseCase.CreateAccount(r.Context(), accountInput)
	if err != nil {
		accController.log.WithError(err).Warn()
	}

	w.Write([]byte(fmt.Sprintf("%s, your account was created", account.Name)))
}
