package main

import (
	"context"
	"log/slog"
	"os"
)

// CustomHandler is a wrapper around slog.Handler that adds the "application" key:value pair to every log entry
type CustomHandler struct {
	handler slog.Handler
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("application", applicationName))
	return h.handler.Handle(ctx, r)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{handler: h.handler.WithGroup(name)}
}

// Define global log handlers
var (
	// stdout including and over level Info
	jsonHandlerStdout = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	customHandlerStdout = &CustomHandler{handler: jsonHandlerStdout}
	logOut              = slog.New(customHandlerStdout)

	// stderr including and over level ERROR
	jsonHandlerStderr = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})
	customHandlerStderr = &CustomHandler{handler: jsonHandlerStderr}
	logErr              = slog.New(customHandlerStderr)
)

// LoggedError outputs the error to stderr as JSON with `"level":"ERROR"`, and returns the error
func LoggedError(err error) error {
	if err != nil {
		logErr.Error(err.Error())
	}
	return err
}
