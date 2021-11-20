package router

import (
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

//GetHTTPHandler returns HTTP handler with all routes
func GetHTTPHandler(dbPool *pgxpool.Pool) *mux.Router {
	accRepo := postgres.NewAccountRepo(dbPool)
	accUseCase := usecase.NewAccountUseCase(accRepo)
	accControler := controller.NewAccountController(accUseCase)

	router := mux.NewRouter()

	//account
	router.HandleFunc("/accounts", accControler.CreateAccount).Methods(http.MethodPost)

	return router
}
