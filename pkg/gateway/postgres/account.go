package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/sqlc"
)

// CreateAccount inserts an account in the database.
func (r Repository) CreateAccount(ctx context.Context, acc entities.Account) error {
	err := sqlc.New(r.conn.GetTxOrPool(ctx)).InsertAccount(ctx, sqlc.InsertAccountParams{
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
				return domain.Error(domain.InvalidParamErrorType, "account already exists", map[string]string{
					"account_id": acc.ID.String(),
				})
			}
		}
		return fmt.Errorf("creating account: %w", err)
	}

	return nil
}

func (r Repository) ListAccounts(ctx context.Context) ([]entities.Account, error) {
	rows, err := sqlc.New(r.conn.GetTxOrPool(ctx)).ListAccounts(ctx)
	if err != nil {
		return nil, domain.NewDBError(domain.QueryRefFetchAcc, err, domain.ErrUnexpected)
	}

	accList := make([]entities.Account, 0, len(rows))
	for _, row := range rows {
		accList = append(accList, parseSqlcAccount(row))
	}

	return accList, nil
}

func (r Repository) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	row, err := sqlc.New(r.conn.GetTxOrPool(ctx)).GetAccount(ctx, uuid.FromStringOrNil(id.String()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrAccNotFound
		}
		return 0, domain.NewDBError(domain.QueryRefGetAcc, err, domain.ErrUnexpected)
	}

	return int(row.Balance), nil
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func (r Repository) UpdateBalance(ctx context.Context, id vos.UUID, transactionAmount int) error {
	err := sqlc.New(r.conn.GetTxOrPool(ctx)).UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
		Amount: int32(transactionAmount),
		ID:     uuid.FromStringOrNil(id.String()),
	})
	if err != nil {
		return domain.NewDBError(domain.QueryRefUpdateBalance, err, domain.ErrUnexpected)
	}

	return nil
}

func (r Repository) GetAccount(ctx context.Context, cpf vos.CPF) (entities.Account, error) {
	row, err := sqlc.New(r.conn.GetTxOrPool(ctx)).GetAccountByDocument(ctx, cpf.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Account{}, entities.ErrAccNotFound
		}
		return entities.Account{}, domain.NewDBError(domain.QueryRefGetAcc, err, domain.ErrUnexpected)
	}

	return parseSqlcAccount(row), nil
}

func parseSqlcAccount(a sqlc.Account) entities.Account {
	return entities.Account{
		ID:        vos.UUID(a.ID.String()),
		Name:      a.Name,
		CPF:       vos.CPF(a.DocumentNumber),
		Secret:    vos.Secret(a.Secret),
		Balance:   int(a.Balance),
		CreatedAt: a.CreatedAt,
	}
}
