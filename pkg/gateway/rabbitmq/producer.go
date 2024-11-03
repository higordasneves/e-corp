package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wagslane/go-rabbitmq"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
)

type Publisher struct {
	P               *rabbitmq.Publisher
	exchange        string
	accCreationBind string
}

func NewPublisher(ctx context.Context, conn *rabbitmq.Conn, config config.RabbitMQConfig) (Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(config.Exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeDurable,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsConfirm,
	)
	if err != nil {
		return Publisher{}, fmt.Errorf("rabbitmq.NewPublisher: %w", err)
	}

	return Publisher{
		P:               publisher,
		exchange:        config.Exchange,
		accCreationBind: config.Bind,
	}, nil
}

func (p Publisher) publish(ctx context.Context, body any, routingKey string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshiling msg: %w", err)
	}

	err = p.P.PublishWithContext(ctx,
		b,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(p.exchange),
	)
	if err != nil {
		return fmt.Errorf("publising message: %w", err)
	}

	return nil
}

func (p Publisher) NotifyAccountCreation(ctx context.Context, account entities.Account) error {
	msg := struct {
		AccountID string    `json:"account_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		AccountID: account.ID.String(),
		CreatedAt: account.CreatedAt,
	}

	err := p.publish(ctx, msg, p.accCreationBind)
	if err != nil {
		return fmt.Errorf("notifying account creation: %w", err)
	}

	return nil
}
