package logger

import (
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
)

func NewFxLogger() fxevent.Logger {
	return &fxevent.SlogLogger{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     slog.LevelError,
			AddSource: false,
		})),
	}
}
