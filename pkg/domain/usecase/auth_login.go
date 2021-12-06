package usecase

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"time"
)

//LoginInput represents information necessary to access a bank account
type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

type Token string

//Login validates credentials then call the func to create a token
func (authUC authUseCase) Login(ctx context.Context, input *LoginInput) (*Token, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	acc, err := authUC.accountRepo.GetAccount(ctx, input.CPF)
	if err != nil {
		return nil, err
	}

	err = acc.CompareHashSecret(input.Secret)
	if err != nil {
		return nil, domainerr.ErrInvalidPass
	}

	return authUC.createAccToken(acc.ID)
}

//createAccToken generates token for account authorization
func (authUC authUseCase) createAccToken(accID vos.UUID) (*Token, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		Issuer:    "login",
		Subject:   string(accID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(authUC.duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(authUC.secretKey))
	if err != nil {
		return nil, domainerr.ErrUnexpected
	}
	accToken := Token(ss)
	return &accToken, nil
}
