package repository

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
)

type AccountRepo interface {
	CreateAccount(context.Context, *models.Account) error
	FetchAccounts(ctx context.Context) ([]models.AccountOutput, error)
}
