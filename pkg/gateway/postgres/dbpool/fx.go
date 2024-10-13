package dbpool

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
)

var Module = fx.Module("dbpool",
	fx.Provide(
		fx.Annotate(
			func(ctx context.Context, lc fx.Lifecycle, cfg config.Config) (*pgxpool.Pool, error) {
				dbPool, err := pgxpool.New(ctx, cfg.DB.DNS())
				if err != nil {
					return nil, fmt.Errorf("creating new pgx pool: %w", err)
				}

				if err = dbPool.Ping(ctx); err != nil {
					return nil, fmt.Errorf("pinging pool: %w", err)
				}

				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						dbPool.Close()
						return nil
					},
				})

				return dbPool, nil
			},
		),
		NewConn,
	),
)
