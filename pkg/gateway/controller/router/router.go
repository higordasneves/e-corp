package router

import (
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/gateway/controller/middleware"
)

type API interface {
	Login(w http.ResponseWriter, r *http.Request)

	GetBalance(w http.ResponseWriter, r *http.Request)
	CreateAccount(w http.ResponseWriter, r *http.Request)
	ListAccounts(w http.ResponseWriter, r *http.Request)

	ListTransfers(w http.ResponseWriter, r *http.Request)
	Transfer(w http.ResponseWriter, r *http.Request)
}

// HTTPHandler returns HTTP handler with all routes.
func HTTPHandler(l *zap.Logger, api API, cfg config.Config) http.Handler {
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.LoggerToContext(l))

	chiRouter.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	apiVersion := "/api/v1"
	chiRouter.Route(apiVersion, func(r chi.Router) {
		// login
		r.Post("/login", api.Login)

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", api.CreateAccount)
			r.Get("/", api.ListAccounts)
			r.Get("/{account_id}/balance", api.GetBalance)
		})

		// transfers
		r.Route("/transfers", func(r chi.Router) {
			r.Use(middleware.Authenticate(cfg.Auth.SecretKey))
			r.Post("/", api.Transfer)
			r.Get("/", api.ListTransfers)
		})

	})

	return chiRouter
}
