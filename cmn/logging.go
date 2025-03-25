package cmn

import (
	"context"
	"log/slog"
	"os"
)

// CustomHandler is a wrapper around slog.Handler that adds the "application" key:value pair to every log entry
type CustomHandler struct {
	ApplicationName string
	Handler         slog.Handler
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("application", h.ApplicationName))
	return h.Handler.Handle(ctx, r)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{Handler: h.Handler.WithGroup(name)}
}

// Define global log handlers
var (
	// stdout including and over level Info
	jsonHandlerStdout = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	customHandlerStdout = &CustomHandler{Handler: jsonHandlerStdout}
	LogOut              = slog.New(customHandlerStdout)

	// stderr including and over level ERROR
	jsonHandlerStderr = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})
	customHandlerStderr = &CustomHandler{Handler: jsonHandlerStderr}
	LogErr              = slog.New(customHandlerStderr)
)

// Allows the calling package to set the application name for the global log handlers
func SetApplicationName(name string) {
	customHandlerStdout.ApplicationName = name
	customHandlerStderr.ApplicationName = name
}

// LoggedError outputs the error to stderr as JSON with `"level":"ERROR"`, and returns the error
func LoggedError(err error) error {
	if err != nil {
		LogErr.Error(err.Error())
	}
	return err
}
