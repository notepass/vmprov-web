package logger

import (
	"log/slog"
	"os"
)

// New creates a new slog logger with JSON output and the given log level.
func New(level string) *slog.Logger {
	lvl, err := parseLevel(level)
	if err != nil {
		lvl = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
}

func parseLevel(level string) (slog.Level, error) {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return slog.LevelInfo, err
	}
	return l, nil
}
