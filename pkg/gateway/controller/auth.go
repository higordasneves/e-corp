package controller

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
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

type LoginRequest struct {
	Document vos.Document `json:"document"`
	Secret   string       `json:"secret"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login validates the credentials of an account and return a login token session.
// It returns bad request error if the password doesn't match.
func (authCtrl AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		reponses.HandleError(w, err, authCtrl.log)
		return
	}

	output, err := authCtrl.authUseCase.Login(r.Context(), usecase.LoginInput(req))
	if err != nil {
		reponses.HandleError(w, err, authCtrl.log)
		return
	}

	// Generating the Claims.
	claims := &jwt.StandardClaims{
		Issuer:    "login",
		Subject:   output.AccountID.String(),
		IssuedAt:  output.IssuedAt.Unix(),
		ExpiresAt: output.ExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(authCtrl.secretKey))
	if err != nil {
		reponses.HandleError(w, err, authCtrl.log)
		return
	}

	reponses.SendResponse(w, http.StatusOK, LoginResponse{Token: tokenString}, authCtrl.log)
}
