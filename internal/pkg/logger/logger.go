package logger

import (
	slog_logger "AggregationService/internal/pkg/logger/slog-logger"
	"context"
	"log/slog"
)

type ctxLogger struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return logger
	}

	return slog_logger.New("")
}
