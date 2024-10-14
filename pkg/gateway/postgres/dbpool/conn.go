package dbpool

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	dbPool *pgxpool.Pool
}

func NewConn(dbPool *pgxpool.Pool) Conn {
	return Conn{dbPool: dbPool}
}

type txCtxKey struct{}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

func (c Conn) BeginTX(ctx context.Context) (context.Context, error) {
	tx, err := c.dbPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning a transaction: %w", err)
	}

	return context.WithValue(ctx, txCtxKey{}, tx), nil
}

func (c Conn) GetTxOrPool(ctx context.Context) Querier {
	if v := ctx.Value(txCtxKey{}); v != nil {
		if pgxTx, ok := v.(*pgxpool.Tx); ok {
			return pgxTx
		}
	}

	return c.dbPool
}

func (c Conn) CommitTX(ctx context.Context) error {
	if v := ctx.Value(txCtxKey{}); v != nil {
		if pgxTx, ok := v.(*pgxpool.Tx); ok {
			err := pgxTx.Commit(ctx)
			if err != nil {
				return fmt.Errorf("committing a transaction: %w", err)
			}

			return nil
		}
	}

	return errors.New("transaction not found")
}

func (c Conn) RollbackTX(ctx context.Context) error {
	if v := ctx.Value(txCtxKey{}); v != nil {
		if pgxTx, ok := v.(*pgxpool.Tx); ok {
			err := pgxTx.Rollback(ctx)
			if err != nil {
				return fmt.Errorf("rolling back a transaction: %w", err)
			}

			return nil
		}
	}

	return errors.New("transaction not found")
}
