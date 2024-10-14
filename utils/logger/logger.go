package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// key is the context key used to associate the logger.
type key struct{}

// New creates a new zap logger.
func New() (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("building logger: %w", err)
	}

	return logger, nil
}

// NewWithCtx creates a new zap logger.
func NewWithCtx(ctx context.Context) (context.Context, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("building logger: %w", err)
	}

	return context.WithValue(ctx, key{}, logger), nil
}

// AssociateCtx returns a new context associated with the given logger.
func AssociateCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

// Logger returns the logger associated with the given context. If there is no logger, it will create a new logger.
func Logger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		panic("nil context passed to logger")
	}

	logger, ok := ctx.Value(key{}).(*zap.Logger)
	if ok {
		return logger
	}

	newLogger, err := New()
	if err != nil {
		panic("creating new logger: " + err.Error())
	}

	return newLogger
}

func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Error(msg, fields...)
}
