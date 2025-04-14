package logging

import (
	"context"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

type loggerKey struct{}

func init() {
	zapLogger, _ := zap.NewProduction()

	// flushes buffer, if any
	defer func() {
		// intentionally ignoring error here, see https://github.com/uber-go/zap/issues/328
		_ = zapLogger.Sync()
	}()

	logger = zapLogger.Sugar()
}

func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey{}).(*zap.SugaredLogger); ok {
		return logger
	}

	return zap.NewNop().Sugar()
}
