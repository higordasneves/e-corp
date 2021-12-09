package controller

import (
	"encoding/json"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
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
		responses.SendError(w, http.StatusBadRequest, err, authCtrl.log)
		return
	}

	var loginInput usecase.LoginInput
	if err = json.Unmarshal(bodyRequest, &loginInput); err != nil {
		responses.SendError(w, http.StatusBadRequest, err, authCtrl.log)
		return
	}

	token, err := authCtrl.authUseCase.Login(r.Context(), &loginInput)

	if err != nil {
		if err == entities.ErrAccNotFound || err == vos.ErrInvalidPass {
			responses.SendError(w, http.StatusBadRequest, err, authCtrl.log)
			return
		}
		responses.SendError(w, http.StatusInternalServerError, err, authCtrl.log)
		return
	}
	responses.SendResponse(w, http.StatusOK, token, authCtrl.log)
}
