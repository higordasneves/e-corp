package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type account struct {
	dbPool *pgxpool.Pool
}

func NewAccountRepo(dbPool *pgxpool.Pool) repository.AccountRepo {
	return &account{dbPool}
}

//CreateAccount inserts a account in database
func (accRepo account) CreateAccount(ctx context.Context, acc *models.Account) error {
	statement, err := accRepo.dbPool.Query(ctx, "INSERT INTO accounts "+
		"(id, cpf, name, secret, balance, created_at)"+
		" VALUES ($1, $2, $3, $4, $5, $6)", acc.ID, acc.CPF, acc.Name, acc.Secret, acc.Balance, acc.CreatedAt)
	defer statement.Close()

	if err != nil {
		return err
	}
	return nil
}
