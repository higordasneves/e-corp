package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/middleware"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/dbpool"
)

// HTTPHandler returns HTTP handler with all routes.
func HTTPHandler(dbPool *pgxpool.Pool, l *zap.Logger, cfgAuth *config.AuthConfig) http.Handler {
	r := postgres.NewRepository(dbpool.NewConn(dbPool))

	accUseCase := usecase.NewAccountUseCase(r)
	accController := controller.NewAccountController(accUseCase)

	tUseCase := usecase.NewTransferUseCase(r)
	tController := controller.NewTransferController(tUseCase)

	authUseCase := usecase.NewAuthUseCase(r, cfgAuth)
	authController := controller.NewAuthController(authUseCase, cfgAuth.SecretKey)

	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.LoggerToContext(l))

	chiRouter.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	apiVersion := "/api/v0"
	chiRouter.Route(apiVersion, func(r chi.Router) {
		// login
		r.Post("/login", authController.Login)

		// accounts
		r.Route("/accounts", func(r chi.Router) {
			r.Post("", accController.CreateAccount)
			r.Get("", accController.CreateAccount)
			r.Get("/{account_id}/balance", accController.GetBalance)
		})

		// transfers
		r.Route("/transfers", func(r chi.Router) {
			r.Use(middleware.Authenticate(cfgAuth.SecretKey))
			r.Post("", tController.Transfer)
			r.Get("", tController.ListTransfers)
		})
	})

	return chiRouter
}
