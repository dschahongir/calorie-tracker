// pkg/logger/logger.go
package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

// Init инициализирует глобальный structured JSON logger
func Init(env string) *slog.Logger {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Для локальной разработки — понятный текстовый формат
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}

// WithContext добавляет логгер в контекст
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext извлекает логгер из контекста
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
