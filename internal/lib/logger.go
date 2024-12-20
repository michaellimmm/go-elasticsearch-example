package lib

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Enables caller information
	}))
}
