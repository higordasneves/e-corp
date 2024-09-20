package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/dbpool"
)

type Repository struct {
	conn dbpool.Conn
}

func NewRepository(dbPool dbpool.Conn) Repository {
	return Repository{dbPool}
}

func (r Repository) BeginTX(ctx context.Context) (context.Context, error) {
	return r.conn.BeginTX(ctx)
}

func (r Repository) CommitTX(ctx context.Context) error {
	return r.conn.CommitTX(ctx)
}

func (r Repository) RollbackTX(ctx context.Context) error {
	return r.conn.RollbackTX(ctx)
}
