package postgres

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PerformTransaction(ctx context.Context, ctxChan chan context.Context, db *pgxpool.Pool, errChan chan error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		close(ctxChan)
		return domain.NewDBError(domain.QueryRefPerformTransaction, fmt.Errorf("cannot begin transaction, %s", err), domain.ErrUnexpected)
	}

	dbCtx := context.WithValue(ctx, "dbConnection", tx)

	ctxChan <- dbCtx
	err = <-errChan

	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return domain.NewDBError(domain.QueryRefPerformTransaction, fmt.Errorf("%s, rollback failed too: %s", err, errRB), domain.ErrUnexpected)
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.NewDBError(domain.QueryRefPerformTransaction, fmt.Errorf("cannot commit transaction, %s", err), domain.ErrUnexpected)
	}
	return nil
}
