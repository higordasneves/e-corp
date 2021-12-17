package main

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/http/router"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg := config.Config{}
	cfg.LoadEnv()

	dbDNS := cfg.DB.DNS()
	log.Info("Accessing database")
	ctxDB := context.Background()
	dbPool, err := pgxpool.Connect(ctxDB, dbDNS)

	defer dbPool.Close()

	if err != nil {
		log.WithError(err).Fatal(config.ErrConnectDB)
	}

	if err = dbPool.Ping(ctxDB); err != nil {
		log.WithError(err).Fatal(config.ErrConnectDB)
	}

	migrationPath := "pkg/gateway/postgres/migrations"
	err = postgres.Migration(migrationPath, dbPool, log)
	if err != nil {
		log.WithError(err).Fatal(config.ErrMigrateDB)
	}

	r := router.GetHTTPHandler(dbPool, log, &cfg.Auth)
	log.Fatal(http.ListenAndServe(":5000", r))
}
