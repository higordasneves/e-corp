package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/fx"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/server"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/dbpool"
	"github.com/higordasneves/e-corp/pkg/gateway/rabbitmq"
	"github.com/higordasneves/e-corp/utils/apictx"
	"github.com/higordasneves/e-corp/utils/logger"
)

// @Title Ecorp API
// @Version 1.0
// @Description A MVP of an API for banking accounts

// @in header
// @name Authorization
func main() {
	app := fx.New(Options)
	if err := app.Err(); err != nil {
		panic(err)
	}

	app.Run()
}

var Options = fx.Options(
	logger.Module,
	apictx.Module,
	config.Module,
	dbpool.Module,
	server.Module,
	rabbitmq.ModuleConn,
	rabbitmq.ModulePub,
	fx.Invoke(func(ctx context.Context, pool *pgxpool.Pool) error {
		err := postgres.Migration(ctx, "pkg/gateway/postgres/migrations", pool)
		if err != nil {
			return fmt.Errorf("executing migrations: %w", err)
		}

		return nil
	}),
	fx.Provide(
		postgres.NewRepository,
		fx.Annotate(
			controller.NewApi,
			fx.As(new(server.API)),
		),
	),
)
