package postgres

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/sqlc"
	"github.com/higordasneves/e-corp/pkg/repository"
)

type account struct {
	dbPool *pgxpool.Pool
}

func NewAccountRepo(dbPool *pgxpool.Pool) repository.AccountRepo {
	return &account{dbPool}
}

// CreateAccount inserts account in database
func (accRepo account) CreateAccount(ctx context.Context, acc *entities.Account) error {
	err := sqlc.New(accRepo.dbPool).InsertAccount(ctx, sqlc.InsertAccountParams{
		ID:             uuid.FromStringOrNil(acc.ID.String()),
		DocumentNumber: acc.CPF.String(),
		Name:           acc.Name,
		Secret:         acc.Secret.String(),
		Balance:        int64(acc.Balance),
		CreatedAt:      acc.CreatedAt,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return entities.ErrAccAlreadyExists
			}
		}
		return repository.NewDBError(repository.QueryRefCreateAcc, err, repository.ErrUnexpected)
	}

	return nil
}

func (accRepo account) FetchAccounts(ctx context.Context) ([]entities.Account, error) {
	rows, err := sqlc.New(accRepo.dbPool).ListAccounts(ctx)
	if err != nil {
		return nil, repository.NewDBError(repository.QueryRefFetchAcc, err, repository.ErrUnexpected)
	}

	accList := make([]entities.Account, 0, len(rows))
	for _, row := range rows {
		accList = append(accList, entities.Account{
			ID:        vos.UUID(row.ID.String()),
			Name:      row.Name,
			CPF:       vos.CPF(row.DocumentNumber),
			Secret:    vos.Secret(row.Secret),
			Balance:   int(row.Balance),
			CreatedAt: row.CreatedAt,
		})
	}

	return accList, nil
}

func (accRepo account) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	row, err := sqlc.New(accRepo.dbPool).GetAccount(ctx, uuid.FromStringOrNil(id.String()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrAccNotFound
		}
		return 0, repository.NewDBError(repository.QueryRefGetAcc, err, repository.ErrUnexpected)
	}

	return int(row.Balance), nil
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func (accRepo account) UpdateBalance(ctx context.Context, id vos.UUID, transactionAmount int) error {
	var db Querier
	db = accRepo.dbPool

	if tx := ctx.Value("dbConnection"); tx != nil {
		db = tx.(*pgxpool.Tx)
	}

	err := sqlc.New(db).UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
		Amount: int32(transactionAmount),
		ID:     uuid.FromStringOrNil(id.String()),
	})
	if err != nil {
		return repository.NewDBError(repository.QueryRefUpdateBalance, err, repository.ErrUnexpected)
	}

	return nil
}

func (accRepo account) GetAccount(ctx context.Context, cpf vos.CPF) (*entities.Account, error) {
	row, err := sqlc.New(accRepo.dbPool).GetAccountByDocument(ctx, cpf.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrAccNotFound
		}
		return nil, repository.NewDBError(repository.QueryRefGetAcc, err, repository.ErrUnexpected)
	}

	return &entities.Account{
		ID:        vos.UUID(row.ID.String()),
		Name:      row.Name,
		CPF:       vos.CPF(row.DocumentNumber),
		Secret:    vos.Secret(row.Secret),
		Balance:   int(row.Balance),
		CreatedAt: row.CreatedAt,
	}, nil
}
