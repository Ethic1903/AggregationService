package sloglogger

import (
	"log"
	"log/slog"
	"os"
)

const (
	_envDev  = "dev"
	_envProd = "prod"
)

func New(env string) *slog.Logger {
	var opts *slog.HandlerOptions

	switch env {
	case _envDev:
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	case _envProd:
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	default:
		log.Println("slog-logger.New: unknown env, fallback to dev")
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	handler := slog.NewTextHandler(
		os.Stdout,
		opts,
	)

	return slog.New(handler)
}
