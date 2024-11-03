package rabbitmq

import (
	"context"
	"fmt"

	"github.com/wagslane/go-rabbitmq"
)

func NewConn(ctx context.Context, url string) (*rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(
		url,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq.NewConn: %w", err)
	}

	return conn, nil
}
