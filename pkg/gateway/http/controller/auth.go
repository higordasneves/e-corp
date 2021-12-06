package controller

import (
	"encoding/json"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/responses"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	authUseCase usecase.AuthUseCase
	log         *logrus.Logger
}

func NewAuthController(authUseCase usecase.AuthUseCase, log *logrus.Logger) AuthController {
	return &authController{authUseCase: authUseCase, log: log}
}

func (authCtrl authController) Login(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendResponse(w, http.StatusBadRequest, responses.ErrorJSON(err), authCtrl.log)
		return
	}

	var loginInput usecase.LoginInput
	if err = json.Unmarshal(bodyRequest, &loginInput); err != nil {
		responses.SendResponse(w, http.StatusBadRequest, responses.ErrorJSON(err), authCtrl.log)
		return
	}

	token, err := authCtrl.authUseCase.Login(r.Context(), &loginInput)

	if err != nil {
		if err == domainerr.ErrAccNotFound || err == domainerr.ErrInvalidPass {
			responses.SendResponse(w, http.StatusBadRequest, responses.ErrorJSON(err), authCtrl.log)
			return
		}
		responses.SendResponse(w, http.StatusInternalServerError, responses.ErrorJSON(err), authCtrl.log)
		return
	}
	responses.SendResponse(w, http.StatusOK, token, authCtrl.log)
}
