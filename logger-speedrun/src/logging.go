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

// LogError logs the error to stderr
func LoggedError(err error) error {
	jsonHandlerErr := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})
	customHandlerErr := &CustomHandler{handler: jsonHandlerErr}
	logErr := slog.New(customHandlerErr)

	if err != nil {
		logErr.Error(err.Error())
	}
	return err
}
