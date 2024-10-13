package apictx

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/utils/logger"
)

var Module = fx.Module("context",
	fx.Provide(
		fx.Annotate(
			func(lc fx.Lifecycle, l *zap.Logger) context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				ctx = logger.AssociateCtx(ctx, l)

				lc.Append(fx.Hook{
					OnStop: func(_ context.Context) error {
						cancel()
						return nil
					},
				})

				return ctx
			},
			fx.As(new(context.Context)),
		),
	),
)
