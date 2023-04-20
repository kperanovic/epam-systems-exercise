// Package logger contains some helper functions
// for managing loggers in go applications.
// The package uses the zap logging library.
package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger context key.
type logKeyType struct{}

var logKey logKeyType

// NewDevelopment creates a new development zap.Logger instance.
func NewDevelopment() *zap.Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return log
}

// NewProduction creates a new production
// zap.Logger instance with a custom timestamp
// field which is compatible with the EFK/ELK stack.
func NewProduction() *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "@timestamp"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.Lock(os.Stderr),
		zap.NewAtomicLevel(),
	))

	return log
}

// WithContext will return a new context
// which contains an instance of zap.Logger.
func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey, log)
}

// FromContext returns a zap.Logger instance
// from the received context if if exists.
// Otherwise the function returns nil.
func FromContext(ctx context.Context) *zap.Logger {
	log := ctx.Value(logKey)
	if log == nil {
		return nil
	}
	return log.(*zap.Logger)
}
