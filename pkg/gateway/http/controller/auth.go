package controller

import (
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"github.com/sirupsen/logrus"
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
	var loginInput usecase.LoginInput
	if err := interpreter.ReadRequestBody(r, &loginInput); err != nil {
		interpreter.HandleError(w, err, authCtrl.log)
		return
	}

	token, err := authCtrl.authUseCase.Login(r.Context(), &loginInput)

	if err != nil {
		interpreter.HandleError(w, err, authCtrl.log)
		return
	}
	interpreter.SendResponse(w, http.StatusOK, token, authCtrl.log)
}
