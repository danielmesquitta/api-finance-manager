package jsonlog

import (
	"log/slog"
	"os"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
)

func NewLogger(e *config.Env) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch e.Environment {
	case config.EnvironmentDevelopment:
		opts.Level = slog.LevelDebug

	default:
		opts.Level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
