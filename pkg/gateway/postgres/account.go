package postgres

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type account struct {
	dbPool *pgxpool.Pool
}

func NewAccountRepo(dbPool *pgxpool.Pool) repository.AccountRepo {
	return &account{dbPool}
}

//CreateAccount inserts account in database
func (accRepo account) CreateAccount(ctx context.Context, acc *entities.Account) error {
	_, err := accRepo.dbPool.Exec(ctx, "INSERT INTO accounts "+
		"(id, cpf, name, secret, balance, created_at)"+
		" VALUES ($1, $2, $3, $4, $5, $6)", acc.ID.String(), acc.CPF, acc.Name, acc.Secret, int64(acc.Balance), acc.CreatedAt)

	var pgErr *pgconn.PgError

	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entities.ErrAccAlreadyExists
			}
		}
		return repository.NewDBError(repository.QueryRefCreateAcc, err)
	}
	return nil
}

func (accRepo account) FetchAccounts(ctx context.Context) ([]entities.Account, error) {
	accCount := accRepo.dbPool.QueryRow(ctx, "select count(*) as count from accounts")

	var count int
	err := accCount.Scan(&count)
	if err != nil {
		return nil, err
	}
	accList := make([]entities.Account, 0, count)

	rows, err := accRepo.dbPool.Query(ctx, "select id, name, cpf, balance::numeric as balance, created_at from accounts")

	defer rows.Close()
	if err != nil {
		return nil, repository.NewDBError(repository.QueryRefFetchAcc, err)
	}

	for rows.Next() {
		var acc entities.Account
		err = rows.Scan(&acc.ID, &acc.Name, &acc.CPF, &acc.Balance, &acc.CreatedAt)
		if err != nil {
			return nil, repository.NewDBError(repository.QueryRefFetchAcc, err)
		}
		accList = append(accList, acc)
	}
	return accList, nil
}

func (accRepo account) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	row := accRepo.dbPool.QueryRow(ctx,
		`select balance
			from accounts
			where id = $1`, id.String())

	var balance int
	err := row.Scan(&balance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, entities.ErrAccNotFound
		}
		return 0, repository.NewDBError(repository.QueryRefGetBalance, err)
	}

	return balance, nil
}

func (accRepo account) UpdateBalance(ctx context.Context, id vos.UUID) (int, error) {
	row := accRepo.dbPool.QueryRow(ctx,
		`select balance
			from accounts
			where id = $1`, id.String())

	var balance int
	err := row.Scan(&balance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, entities.ErrAccNotFound
		}
		return 0, repository.NewDBError(repository.QueryRefUpdateBalance, err)
	}

	return balance, nil
}

func (accRepo account) GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	row := accRepo.dbPool.QueryRow(ctx,
		`select id
			, name
			, cpf
			, secret
			, balance
			, created_at
			from accounts
			where cpf = $1`, cpf)

	var acc entities.Account
	err := row.Scan(&acc.ID, &acc.Name, &acc.CPF, &acc.Secret, &acc.Balance, &acc.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrAccNotFound
		}

		return nil, repository.NewDBError(repository.QueryRefGetAcc, err)
	}

	return &acc, nil
}
