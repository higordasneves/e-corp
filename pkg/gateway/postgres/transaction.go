package postgres

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

func PerformTransaction(ctx context.Context, ctxChan chan context.Context, db *pgxpool.Pool, errChan chan error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		close(ctxChan)
		return repository.NewDBError(repository.QueryRefTransfer, fmt.Errorf("cannot begin transaction, %s", err))
	}

	dbCtx := context.WithValue(ctx, "dbConnection", tx)

	ctxChan <- dbCtx
	err = <-errChan

	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return repository.NewDBError(repository.QueryRefTransfer, fmt.Errorf("%s, rollback failed too: %s", err, errRB))
		}
		return repository.NewDBError(repository.QueryRefTransfer, fmt.Errorf("%s, rollback was performed", err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return repository.NewDBError(repository.QueryRefTransfer, fmt.Errorf("cannot commit transaction, %s", err))
	}
	return nil
}
