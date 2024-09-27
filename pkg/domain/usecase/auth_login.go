package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

// LoginInput represents information necessary to access a bank account
type LoginInput struct {
	CPF    vos.Document `json:"cpf"`
	Secret string       `json:"secret"`
}

type LoginToken string

// Login validates credentials then call the func to generate a login session token with expiration.
func (authUC AuthUseCase) Login(ctx context.Context, input LoginInput) (LoginToken, error) {
	acc, err := authUC.accountRepo.GetAccountByDocument(ctx, input.CPF)
	if err != nil {
		return "", err
	}

	err = acc.Secret.CompareHashSecret(input.Secret)
	if err != nil {
		return "", fmt.Errorf("%w: %w", domain.ErrInvalidParameter, err)
	}

	token, err := authUC.generateToken(acc.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

// generateToken generates token for account authorization.
func (authUC AuthUseCase) generateToken(accID uuid.UUID) (LoginToken, error) {
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
		return "", err
	}

	return LoginToken(ss), nil
}
