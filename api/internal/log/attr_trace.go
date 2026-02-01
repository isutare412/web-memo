package log

import (
	"context"
	"log/slog"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/otel/trace"
)

// newAttrTraceMiddleware creates a middleware that adds trace ID from the
// tracing span in the context to the log record.
func newAttrTraceMiddleware() slogmulti.Middleware {
	return slogmulti.NewInlineMiddleware(
		func(ctx context.Context, level slog.Level, next func(context.Context, slog.Level) bool) bool {
			return next(ctx, level)
		},

		func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
			spanCtx := trace.SpanContextFromContext(ctx)
			if spanCtx.IsValid() {
				record.AddAttrs(slog.String("traceId", spanCtx.TraceID().String()))
			}

			return next(ctx, record)
		},

		func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler { return next(attrs) },
		func(name string, next func(string) slog.Handler) slog.Handler { return next(name) },
	)
}
