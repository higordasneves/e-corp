package controller

import (
	"encoding/json"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"io/ioutil"
	"log"
	"net/http"
)

type AccountController interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
}

type accountController struct {
	accUseCase usecase.AccountUseCase
}

func NewAccountController(accUseCase usecase.AccountUseCase) AccountController {
	return &accountController{accUseCase: accUseCase}
}

func (accController accountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var accountInput usecase.AccountInput
	if err := json.Unmarshal(bodyRequest, &accountInput); err != nil {
		log.Print(err)
	}

	account, err := accController.accUseCase.CreateAccount(accountInput)
	if err != nil {
		log.Print(err.Error())
	}

	w.Write([]byte(fmt.Sprintf("%s, your account was created", account.Name)))
}
