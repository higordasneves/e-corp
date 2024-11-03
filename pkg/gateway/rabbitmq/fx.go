package rabbitmq

import (
	"context"
	"fmt"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/utils/logger"
)

var ModuleConn = fx.Module("rabbitmq-conn",
	fx.Provide(
		fx.Annotate(
			func(ctx context.Context, lc fx.Lifecycle, cfg config.Config) (*rabbitmq.Conn, error) {
				conn, err := NewConn(ctx, cfg.MQ.URL())
				if err != nil {
					return nil, fmt.Errorf("creating rabbit conn: %w", err)
				}

				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						conn.Close()
						return nil
					},
				})

				return conn, nil
			},
		),
	),
)

var ModulePub = fx.Module("rabbitmq-pub",
	fx.Provide(
		fx.Annotate(
			func(ctx context.Context, lc fx.Lifecycle, cfg config.Config, conn *rabbitmq.Conn) (Publisher, error) {
				pub, err := NewPublisher(ctx, conn, cfg.MQ)
				if err != nil {
					return Publisher{}, fmt.Errorf("creating rabbit publisher: %w", err)
				}

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						pub.P.NotifyReturn(func(r rabbitmq.Return) {
							logger.Error(ctx, "rabbitmq notify received",
								zap.String("message_id", r.MessageId),
								zap.String("body", string(r.Body)),
							)
						})

						pub.P.NotifyPublish(func(c rabbitmq.Confirmation) {
							logger.Info(ctx, "rabbitmq confirmation received",
								zap.Any("confirmation", c.Confirmation),
							)
						})
						return nil
					},
				})

				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						pub.P.Close()
						return nil
					},
				})

				return pub, nil
			},
		),
	),
)

var ModuleSub = fx.Module("rabbitmq-sub",
	fx.Provide(
		fx.Annotate(
			func(ctx context.Context, lc fx.Lifecycle, cfg config.Config, conn *rabbitmq.Conn) (Consumer, error) {
				consumer, err := NewConsumer(ctx, conn, cfg.MQ)
				if err != nil {
					return Consumer{}, fmt.Errorf("creating rabbit consumer: %w", err)
				}

				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						consumer.C.Close()
						return nil
					},
				})

				return consumer, nil
			},
		),
	),
)
