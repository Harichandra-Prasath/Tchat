package logging

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func IntialiseLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	Logger = slog.New(handler)
}
