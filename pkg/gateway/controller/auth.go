package controller

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
)

//go:generate moq -stub -pkg mocks -out mocks/auth_uc.go . AuthUseCase

type AuthUseCase interface {
	Login(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error)
}

type AuthController struct {
	authUseCase AuthUseCase
	secretKey   string
}

func NewAuthController(authUseCase AuthUseCase, secretKey string) AuthController {
	return AuthController{authUseCase, secretKey}
}

type LoginRequest struct {
	Document vos.Document `json:"document"`
	Secret   string       `json:"secret"`
}

type LoginResponse struct {
	// Token is the session token used to authenticate the account.
	Token string `json:"token"`
}

// Login validates the credentials of an account and return a login token session.
// It returns bad request error if the password doesn't match.
// @Summary Login
// @Description Validates the credentials of an account and return a login token session.
// @Description It returns bad request error if the provided password doesn't match for the account.
// @Tags Login
// @Param Body body LoginRequest true "Request body"
// @Accept json
// @Produce json
// @Success 200 {object} LoginResponse "Token"
// @Failure 400 {object} ErrorResponse "invalid parameter"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/login [POST]
func (authCtrl AuthController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		HandleError(ctx, w, err)
		return
	}

	output, err := authCtrl.authUseCase.Login(r.Context(), usecase.LoginInput(req))
	if err != nil {
		HandleError(ctx, w, err)
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
		HandleError(ctx, w, err)
		return
	}

	SendResponse(ctx, w, http.StatusOK, LoginResponse{Token: tokenString})
}
