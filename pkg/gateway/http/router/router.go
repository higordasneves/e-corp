package router

import (
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/http/middleware"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

//GetHTTPHandler returns HTTP handler with all routes
func GetHTTPHandler(dbPool *pgxpool.Pool, log *logrus.Logger, cfgAuth *config.AuthConfig) *mux.Router {
	accRepo := postgres.NewAccountRepo(dbPool)
	accUseCase := usecase.NewAccountUseCase(accRepo)
	accController := controller.NewAccountController(accUseCase, log)

	tRepo := postgres.NewTransferRepository(dbPool)
	tUseCase := usecase.NewTransferUseCase(accRepo, tRepo)
	tController := controller.NewTransferController(tUseCase, log)

	authUseCase := usecase.NewAuthUseCase(accRepo, cfgAuth)
	authController := controller.NewAuthController(authUseCase, log)

	router := mux.NewRouter()

	//account
	router.HandleFunc("/accounts", accController.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts", accController.FetchAccounts).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/balance", accController.GetBalance).Methods(http.MethodGet)

	//transfer
	router.HandleFunc("/transfers", middleware.Authenticate(authUseCase, tController.Transfer, log)).Methods(http.MethodPost)
	router.HandleFunc("/transfers", middleware.Authenticate(authUseCase, tController.FetchTransfers, log)).Methods(http.MethodGet)

	//login
	router.HandleFunc("/login", authController.Login).Methods(http.MethodPost)

	return router
}
