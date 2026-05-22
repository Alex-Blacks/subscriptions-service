package logging

import (
	"context"
	"log/slog"
	"os"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func NewLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
}

func LoggerWithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
