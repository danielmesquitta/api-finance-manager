package config

import (
	"log/slog"
	"os"
)

func init() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(slogger)
}
