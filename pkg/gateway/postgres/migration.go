package postgres

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func Migration(migrationPath string, dbPool *pgxpool.Pool, log *logrus.Logger) error {

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

	if err == migrate.ErrNoChange {
		log.WithError(err).Warn(config.ErrMigrateDB)
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
