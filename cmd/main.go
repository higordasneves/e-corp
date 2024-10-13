package main

import (
	"context"
	"go.uber.org/fx"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/server"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/dbpool"
	"github.com/higordasneves/e-corp/utils/apictx"
	"github.com/higordasneves/e-corp/utils/logger"
)

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
	fx.Invoke(func(ctx context.Context, pool *pgxpool.Pool) {
		postgres.Migration(ctx, "pkg/gateway/postgres/migrations", pool)

	}),
	fx.Provide(
		postgres.NewRepository,
		fx.Annotate(
			newAPI,
			fx.As(new(server.API)),
		),
	),
)

func newAPI(r postgres.Repository, cfg config.Config) controller.API {
	accUseCase := usecase.NewAccountUseCase(r)
	accController := controller.NewAccountController(accUseCase)

	tUseCase := usecase.NewTransferUseCase(r)
	tController := controller.NewTransferController(tUseCase)

	authUseCase := usecase.NewAuthUseCase(r, &cfg.Auth)
	authController := controller.NewAuthController(authUseCase, cfg.Auth.SecretKey)

	return controller.API{
		AuthController:     authController,
		AccountController:  accController,
		TransferController: tController,
	}
}
