package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey int

const (
	requestIDKey contextKey = iota
	traceIDKey
)

func New(level, format string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}
	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func TraceIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(traceIDKey).(string)
	return v
}

// FromContext returns base logger with request_id and trace_id attached from ctx.
func FromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
	args := make([]any, 0, 4)
	if rid := RequestIDFromContext(ctx); rid != "" {
		args = append(args, "request_id", rid)
	}
	if tid := TraceIDFromContext(ctx); tid != "" {
		args = append(args, "trace_id", tid)
	}
	if len(args) == 0 {
		return base
	}
	return base.With(args...)
}
