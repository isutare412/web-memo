package log

import (
	"context"
	"log/slog"

	slogmulti "github.com/samber/slog-multi"
)

// newAttrConstantMiddleware creates a middleware that adds constant key-value
// attributes to all log records.
func newAttrConstantMiddleware(attrs ...slog.Attr) slogmulti.Middleware {
	return slogmulti.NewHandleInlineMiddleware(
		func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
			record.AddAttrs(attrs...)
			return next(ctx, record)
		},
	)
}
