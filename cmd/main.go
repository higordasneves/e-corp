package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/router"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/utils/logger"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic("creating logger: " + err.Error())
	}
	ctx := logger.AssociateCtx(context.Background(), log)

	cfg := config.Config{}
	cfg.LoadEnv()

	dbDNS := cfg.DB.DNS()
	log.Info("Accessing database")
	ctxDB := context.Background()
	dbPool, err := pgxpool.New(ctxDB, dbDNS)
	if err != nil {
		log.Error("creating new pgx pool", zap.Error(err))
	}
	defer dbPool.Close()

	if err = dbPool.Ping(ctxDB); err != nil {
		log.Error("connecting the database", zap.Error(err))
	}

	migrationPath := "pkg/gateway/postgres/migrations"
	err = postgres.Migration(ctx, migrationPath, dbPool)
	if err != nil {
		log.Error("executing database migration", zap.Error(err))
	}

	r := router.HTTPHandler(dbPool, log, &cfg.Auth)
	if err := http.ListenAndServe(":5000", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("failed to start gateway HTTP server", zap.Error(err))
	}
}
