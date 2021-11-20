package postgres

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

func Migration(dbPool *pgxpool.Pool, log *logrus.Logger) {

	cfg := dbPool.Config().ConnConfig

	db := stdlib.OpenDB(*cfg)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.WithError(err).Fatal(config.ErrConnectDB)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/gateway/postgres/migrations",
		"postgres", driver)
	if err != nil {
		log.WithError(err).Fatal(config.ErrMigrateDB)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.WithError(err).Warn(config.ErrMigrateDB)
	} else if err != nil {
		log.WithError(err).Fatal(config.ErrMigrateDB)
	}
}
