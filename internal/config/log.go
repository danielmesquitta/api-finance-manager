package config

import (
	"log/slog"
	"os"
)

func setDefaultLogger(e *Env) {
	opts := &slog.HandlerOptions{}

	switch e.Environment {
	case EnvironmentDevelopment:
		opts.Level = slog.LevelDebug
	case EnvironmentTest:
		opts.Level = slog.LevelDebug
	default:
		opts.Level = slog.LevelInfo
	}

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(slogger)
}
