package router

import (
	"net/http"

	"github.com/gorilla/mux"
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
func HTTPHandler(dbPool *pgxpool.Pool, l *zap.Logger, cfgAuth *config.AuthConfig) *mux.Router {
	r := postgres.NewRepository(dbpool.NewConn(dbPool))

	accUseCase := usecase.NewAccountUseCase(r)
	accController := controller.NewAccountController(accUseCase)

	//tUseCase := usecase.NewTransferUseCase(r)
	tController := controller.NewTransferController(nil)

	authUseCase := usecase.NewAuthUseCase(r, cfgAuth)
	authController := controller.NewAuthController(authUseCase, cfgAuth.SecretKey)

	router := mux.NewRouter()
	apiVersion := "/api/v0"

	router.Use(middleware.LoggerToContext(l))
	//account
	router.HandleFunc(apiVersion+"/accounts", accController.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc(apiVersion+"/accounts", accController.ListAccounts).Methods(http.MethodGet)
	router.HandleFunc(apiVersion+"/accounts/{account_id}/balance", accController.GetBalance).Methods(http.MethodGet)

	//transfer
	router.HandleFunc(apiVersion+"/transfers", middleware.Authenticate(cfgAuth.SecretKey, tController.Transfer)).Methods(http.MethodPost)
	router.HandleFunc(apiVersion+"/transfers", middleware.Authenticate(cfgAuth.SecretKey, tController.ListTransfers)).Methods(http.MethodGet)

	//login
	router.HandleFunc(apiVersion+"/login", authController.Login).Methods(http.MethodPost)

	return router
}
