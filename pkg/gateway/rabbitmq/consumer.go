package rabbitmq

import (
	"context"
	"fmt"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/utils/logger"
)

type Consumer struct {
	C *rabbitmq.Consumer
}

func NewConsumer(ctx context.Context, conn *rabbitmq.Conn, config config.RabbitMQConfig) (Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		config.Queue,
		rabbitmq.WithConsumerOptionsRoutingKey(config.Bind),
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsExchangeName(config.Exchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
	)
	if err != nil {
		return Consumer{}, fmt.Errorf("creating consumer: %w", err)
	}

	return Consumer{
		C: consumer,
	}, nil
}

func (c Consumer) Run(ctx context.Context) error {
	err := c.C.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		logger.Info(ctx, "rabbitmq msg consumed",
			zap.Any("message_id", d.MessageId),
			zap.String("body", string(d.Body)),
		)

		return rabbitmq.Ack
	})
	if err != nil {
		logger.Error(ctx, "running consumer", zap.Error(err))
	}

	return nil
}
