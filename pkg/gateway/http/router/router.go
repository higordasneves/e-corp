package router

import (
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/http/middleware"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/dbpool"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetHTTPHandler returns HTTP handler with all routes
func GetHTTPHandler(dbPool *pgxpool.Pool, log *logrus.Logger, cfgAuth *config.AuthConfig) *mux.Router {
	r := postgres.NewRepository(dbpool.NewConn(dbPool))

	accUseCase := usecase.NewAccountUseCase(r)
	accController := controller.NewAccountController(accUseCase, log)

	tUseCase := usecase.NewTransferUseCase(r)
	tController := controller.NewTransferController(tUseCase, log)

	authUseCase := usecase.NewAuthUseCase(r, cfgAuth)
	authController := controller.NewAuthController(authUseCase, log)

	router := mux.NewRouter()
	apiVersion := "/api/v0"

	//account
	router.HandleFunc(apiVersion+"/accounts", accController.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc(apiVersion+"/accounts", accController.FetchAccounts).Methods(http.MethodGet)
	router.HandleFunc(apiVersion+"/accounts/{account_id}/balance", accController.GetBalance).Methods(http.MethodGet)

	//transfer
	router.HandleFunc(apiVersion+"/transfers", middleware.Authenticate(authUseCase, tController.Transfer, log)).Methods(http.MethodPost)
	router.HandleFunc(apiVersion+"/transfers", middleware.Authenticate(authUseCase, tController.FetchTransfers, log)).Methods(http.MethodGet)

	//login
	router.HandleFunc(apiVersion+"/login", authController.Login).Methods(http.MethodPost)

	return router
}
