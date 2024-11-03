package main

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/fx"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/rabbitmq"
	"github.com/higordasneves/e-corp/utils/apictx"
	"github.com/higordasneves/e-corp/utils/logger"
)

// @Title Ecorp API
// @Version 1.0
// @Description API for banking accounts

// @in header
// @name Authorization
func main() {
	app := fx.New(Options)
	if err := app.Err(); err != nil {
		panic(err)
	}

	app.Run()
}

var Options = fx.Options(
	logger.Module,
	apictx.Module,
	config.Module,
	rabbitmq.ModuleConn,
	rabbitmq.ModuleSub,
	fx.Invoke(func(ctx context.Context, c rabbitmq.Consumer) error {
		if err := c.Run(ctx); err != nil {
			return fmt.Errorf("failed to start consumer: %w", err)
		}

		return nil
	}),
)
