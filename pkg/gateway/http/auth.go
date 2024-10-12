package http

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/gateway/http/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/http/requests"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
)

//go:generate moq -stub -pkg mocks -out mocks/auth_uc.go . AuthUseCase

type AuthUseCase interface {
	Login(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error)
}

type AuthController struct {
	authUseCase AuthUseCase
	secretKey   string
	log         *logrus.Logger
}

func NewAuthController(authUseCase AuthUseCase, secretKey string, log *logrus.Logger) AuthController {
	return AuthController{authUseCase, secretKey, log}
}

func (authCtrl AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginInput usecase.LoginInput
	if err := requests.ReadRequestBody(r, &loginInput); err != nil {
		reponses.HandleError(w, err, authCtrl.log)
		return
	}

	output, err := authCtrl.authUseCase.Login(r.Context(), loginInput)
	if err != nil {
		reponses.HandleError(w, err, authCtrl.log)
		return
	}

	// Create the Claims
	claims := &jwt.StandardClaims{
		Issuer:    "login",
		Subject:   output.AccountID.String(),
		IssuedAt:  output.IssuedAt.Unix(),
		ExpiresAt: output.ExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	resp, err := token.SignedString([]byte(authCtrl.secretKey))
	if err != nil {
		reponses.HandleError(w, err, authCtrl.log)
		return
	}

	reponses.SendResponse(w, http.StatusOK, resp, authCtrl.log)
}
