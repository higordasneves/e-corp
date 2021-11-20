package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type account struct {
	dbPool *pgxpool.Pool
	log    *logrus.Logger
}

func NewAccountRepo(dbPool *pgxpool.Pool, log *logrus.Logger) repository.AccountRepo {
	return &account{dbPool, log}
}

//CreateAccount inserts a account in database
func (accRepo account) CreateAccount(ctx context.Context, acc *models.Account) error {
	_, err := accRepo.dbPool.Exec(ctx, "INSERT INTO accounts "+
		"(id, cpf, name, secret, balance, created_at)"+
		" VALUES ($1, $2, $3, $4, $5, $6)", acc.ID.String(), acc.CPF, acc.Name, acc.Secret, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
