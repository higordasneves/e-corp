package router

import (
	http2 "github.com/higordasneves/e-corp/pkg/gateway/controller"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/middleware"
)

// GetHTTPHandler returns HTTP handler with all routes
func GetHTTPHandler(dbPool *pgxpool.Pool, log *logrus.Logger, cfgAuth *config.AuthConfig) *mux.Router {
	// r := postgres.NewRepository(dbpool.NewConn(dbPool))

	// accUseCase := usecase.NewAccountUseCase(r)
	accController := http2.NewAccountController(nil, log)

	//tUseCase := usecase.NewTransferUseCase(r)
	tController := http2.NewTransferController(nil, log)

	//authUseCase := usecase.NewAuthUseCase(r, cfgAuth)
	authController := http2.NewAuthController(nil, cfgAuth.SecretKey, log)

	router := mux.NewRouter()
	apiVersion := "/api/v0"

	//account
	router.HandleFunc(apiVersion+"/accounts", accController.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc(apiVersion+"/accounts", accController.FetchAccounts).Methods(http.MethodGet)
	router.HandleFunc(apiVersion+"/accounts/{account_id}/balance", accController.GetBalance).Methods(http.MethodGet)

	//transfer
	router.HandleFunc(apiVersion+"/transfers", middleware.Authenticate(cfgAuth.SecretKey, tController.Transfer, log)).Methods(http.MethodPost)
	router.HandleFunc(apiVersion+"/transfers", middleware.Authenticate(cfgAuth.SecretKey, tController.FetchTransfers, log)).Methods(http.MethodGet)

	//login
	router.HandleFunc(apiVersion+"/login", authController.Login).Methods(http.MethodPost)

	return router
}
