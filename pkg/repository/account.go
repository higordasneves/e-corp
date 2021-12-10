package repository

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

type AccountRepo interface {
	CreateAccount(context.Context, *entities.Account) error
	FetchAccounts(ctx context.Context) ([]entities.Account, error)
	GetBalance(ctx context.Context, id vos.UUID) (int, error)
	GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error)
}
