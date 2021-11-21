package postgres

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	ErrAccNotFound = errors.New("account not found")
)

type account struct {
	dbPool *pgxpool.Pool
	log    *logrus.Logger
}

func NewAccountRepo(dbPool *pgxpool.Pool, log *logrus.Logger) repository.AccountRepo {
	return &account{dbPool, log}
}

//CreateAccount inserts account in database
func (accRepo account) CreateAccount(ctx context.Context, acc *models.Account) error {
	_, err := accRepo.dbPool.Exec(ctx, "INSERT INTO accounts "+
		"(id, cpf, name, secret, balance, created_at)"+
		" VALUES ($1, $2, $3, $4, $5, $6)", acc.ID.String(), acc.CPF, acc.Name, acc.Secret, int64(acc.Balance), acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (accRepo account) FetchAccounts(ctx context.Context) ([]models.AccountOutput, error) {
	accCount := accRepo.dbPool.QueryRow(ctx, "select count(*) as count from accounts")

	var count int
	err := accCount.Scan(&count)
	if err != nil {
		return nil, err
	}
	accList := make([]models.AccountOutput, 0, count)

	rows, err := accRepo.dbPool.Query(ctx, "select id, name, cpf, balance, created_at from accounts")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var acc models.Account
		err = rows.Scan(&acc.ID, &acc.Name, &acc.CPF, &acc.Balance, &acc.CreatedAt)
		accOutput := acc.GetAccOutput()
		if err != nil {
			return nil, err
		}
		accList = append(accList, *accOutput)
	}

	return accList, nil
}

func (accRepo account) GetBalance(ctx context.Context, id vos.AccountID) (*vos.Currency, error) {
	row := accRepo.dbPool.QueryRow(ctx,
		`select balance
			from accounts
			where id = $1`, id.String())

	var balance vos.Currency
	err := row.Scan(&balance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrAccNotFound
		}
		return nil, err
	}

	balance.ConvertFromCents()
	return &balance, nil
}
