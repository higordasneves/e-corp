package router

import (
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

//GetHTTPHandler returns HTTP handler with all routes
func GetHTTPHandler(dbPool *pgxpool.Pool, log *logrus.Logger) *mux.Router {
	accRepo := postgres.NewAccountRepo(dbPool, log)
	accUseCase := usecase.NewAccountUseCase(accRepo, log)
	accControler := controller.NewAccountController(accUseCase, log)

	router := mux.NewRouter()

	//account
	router.HandleFunc("/accounts", accControler.CreateAccount).Methods(http.MethodPost)

	return router
}
