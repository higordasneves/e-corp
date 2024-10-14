package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
)

var Module = fx.Module("httpserver",
	fx.Invoke(
		fx.Annotate(
			func(ctx context.Context, lc fx.Lifecycle, l *zap.Logger, api API, cfg config.Config) (*http.Server, error) {
				handler := HTTPHandler(l, api, cfg)

				server := http.Server{
					Addr:              cfg.HTTP.Address + ":" + cfg.HTTP.Port,
					Handler:           handler,
					ReadTimeout:       time.Second * 30,
					ReadHeaderTimeout: time.Second * 30,
					WriteTimeout:      time.Second * 30,
				}

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							l.Info("HTTP server listening", zap.String("address", server.Addr))

							if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
								l.Error("listening HTTP", zap.Error(err))
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						return server.Shutdown(ctx)
					},
				})

				return &server, nil
			},
		),
	),
)
