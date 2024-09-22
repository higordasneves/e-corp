package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"

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

// ListAccounts Lists accounts by filtering the IDs provided in the input.
// The field LastFetchedID is a cursor and represents the ID of
// the last account listed (on the previous page).
// The query is sorted in descending order.
func (r Repository) ListAccounts(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error) {
	rows, err := sqlc.New(r.conn.GetTxOrPool(ctx)).ListAccounts(ctx, sqlc.ListAccountsParams{
		Ids:           input.IDs,
		LastFetchedID: input.LastFetchedID,
		// We list page size + 1 to check if there will be more items to list on the next page.
		PageSize: int32(input.PageSize) + 1,
	})
	if err != nil {
		return usecase.ListAccountsOutput{}, fmt.Errorf("listing accounts: %w", err)
	}

	var nextPage *usecase.ListAccountsInput
	// If the number of returned items is equal to page size + 1, there will be a next page.
	// We need to construct the cursor.
	if len(rows) >= input.PageSize+1 {
		nextPage = &input
		rows = rows[:len(rows)-1]
		nextPage.LastFetchedID = rows[len(rows)-1].ID
	}

	accList := make([]entities.Account, 0, len(rows))
	for _, row := range rows {
		accList = append(accList, parseSqlcAccount(row))
	}

	return usecase.ListAccountsOutput{
		Accounts: accList,
		NextPage: nextPage,
	}, nil
}

// GetBalance returns the balance of the account for the provided ID.
func (r Repository) GetBalance(ctx context.Context, id vos.UUID) (int, error) {
	row, err := sqlc.New(r.conn.GetTxOrPool(ctx)).GetAccount(ctx, uuid.FromStringOrNil(id.String()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.Error(domain.NotFoundErrorType, "account not found", nil)
		}
		return 0, fmt.Errorf("getting balance: %w", err)
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
