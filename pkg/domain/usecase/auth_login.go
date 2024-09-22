package usecase

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

// LoginInput represents information necessary to access a bank account
type LoginInput struct {
	CPF    vos.CPF `json:"cpf"`
	Secret string  `json:"secret"`
}

type Token string

// Login validates credentials then call the func to create a token
func (authUC AuthUseCase) Login(ctx context.Context, input *LoginInput) (*Token, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	acc, err := authUC.accountRepo.GetAccountByDocument(ctx, input.CPF)
	if err != nil {
		return nil, err
	}

	err = input.CPF.ValidateInput()
	if err != nil {
		return nil, err
	}

	err = acc.Secret.CompareHashSecret(input.Secret)
	if err != nil {
		return nil, vos.ErrInvalidPass
	}

	return authUC.createAccToken(acc.ID)
}

// createAccToken generates token for account authorization
func (authUC AuthUseCase) createAccToken(accID uuid.UUID) (*Token, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		Issuer:    "login",
		Subject:   accID.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(authUC.duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(authUC.secretKey))
	if err != nil {
		return nil, err
	}
	accToken := Token(ss)
	return &accToken, nil
}
