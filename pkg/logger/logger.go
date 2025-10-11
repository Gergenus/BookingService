package logger

import (
	"log/slog"
	"os"
)

func SetUp(logLevel string) *slog.Logger {
	switch logLevel {
	case "debug":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("incorrect loglevel")
	}

}
