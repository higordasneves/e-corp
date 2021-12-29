package controller

import (
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/sirupsen/logrus"
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
