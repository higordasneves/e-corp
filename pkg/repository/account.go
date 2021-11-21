package repository

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountRepo interface {
	CreateAccount(context.Context, *models.Account) error
	FetchAccounts(ctx context.Context) ([]models.AccountOutput, error)
	GetBalance(ctx context.Context, id vos.AccountID) (*vos.Currency, error)
}
