package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/utils/logger"
)

func Migration(ctx context.Context, migrationPath string, dbPool *pgxpool.Pool) error {
	cfg := dbPool.Config().ConnConfig

	db := stdlib.OpenDB(*cfg)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("creating migration: %w", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Error(ctx, "migration error no change", zap.Error(err))
			return nil
		}

		return fmt.Errorf("executin migration: %w", err)
	}

	return nil
}
