package router

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"net/http"
)

func GetHTTPHandler(db *sql.DB) *mux.Router {
	accRepo := postgres.NewAccountRepo(db)
	accUseCase := usecase.NewAccountUseCase(accRepo)
	accControler := controller.NewAccountController(accUseCase)

	router := mux.NewRouter()

	//account
	router.HandleFunc("/account", accControler.CreateAccount).Methods(http.MethodPost)

	return router
}
