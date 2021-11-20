package main

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/http/router"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"net/http"
)

func main() {
	log := logrus.New()
	dbCfg := config.DatabaseConfig{}
	dbCfg.LoadEnv()
	dbDNS := dbCfg.DNS()

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

	r := router.GetHTTPHandler(dbPool, log)
	log.Fatal(http.ListenAndServe(":8080", r))

}
