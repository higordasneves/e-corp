package postgres

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/higordasneves/e-corp/utils/logger"
)

func Migration(ctx context.Context, migrationPath string, dbPool *pgxpool.Pool) error {
	cfg := dbPool.Config().ConnConfig

	db := stdlib.OpenDB(*cfg)

	err := db.Ping()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Error(ctx, "migration error no change", zap.Error(err))
			return nil
		}

		return err
	}

	return nil
}
