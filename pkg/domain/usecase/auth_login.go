package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
)

type AuthUCRepository interface {
	GetAccountByDocument(ctx context.Context, cpf vos.Document) (entities.Account, error)
}

type AuthUC struct {
	accountRepo AuthUCRepository
	duration    time.Duration
}

func NewAuthUC(accountRepo AuthUCRepository, cfgAuth *config.AuthConfig) AuthUC {
	return AuthUC{accountRepo: accountRepo, duration: cfgAuth.Duration}
}

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
func (uc AuthUC) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
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
