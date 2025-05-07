package logger

import (
	"log/slog"
	"os"
)

// New creates a new logger instance with default options
func New() *slog.Logger {
	// Create a new logger with TextHandler and the specified options
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}

// NewWithLevel creates a new logger instance with the specified log level
func NewWithLevel(level slog.Level) *slog.Logger {
	// Set handler options to enable logging at the specified level and above
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Create a new logger with TextHandler and the specified options
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return logger
}
