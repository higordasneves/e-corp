package controller

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
)

//go:generate moq -stub -pkg mocks -out mocks/auth_uc.go . AuthUseCase

type AuthUseCase interface {
	Login(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error)
	ValidateToken(tokenString string) (*jwt.StandardClaims, error)
}

type AuthController struct {
	authUseCase AuthUseCase
	log         *logrus.Logger
}

func NewAuthController(authUseCase AuthUseCase, log *logrus.Logger) AuthController {
	return AuthController{authUseCase: authUseCase, log: log}
}

func (authCtrl AuthController) Login(w http.ResponseWriter, r *http.Request) {
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
