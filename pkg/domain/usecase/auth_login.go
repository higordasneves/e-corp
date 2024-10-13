package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

// LoginInput represents information necessary to access a bank account.
type LoginInput struct {
	Document vos.Document
	Secret   string
}

type LoginOutput struct {
	AccountID uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type LoginToken string

// Login validates the credentials of an account and return a login token session.
// It returns domain.ErrInvalidParameter if the password doesn't match.
func (uc AuthUseCase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	acc, err := uc.accountRepo.GetAccountByDocument(ctx, input.Document)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("getting account: %w", err)
	}

	err = acc.Secret.CompareHashSecret(input.Secret)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("%w: %w", domain.ErrInvalidParameter, err)
	}

	now := time.Now()

	return LoginOutput{
		AccountID: acc.ID,
		IssuedAt:  now,
		ExpiresAt: now.Add(uc.duration),
	}, nil
}
