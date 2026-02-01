package log

import (
	"context"
	"log/slog"

	slogmulti "github.com/samber/slog-multi"
)

type attrContextKey struct{}

// attrContext holds attributes that can be added to log records.
type attrContext struct {
	attrs []slog.Attr
}

// WithAttrContext creates a new context with an empty attribute context. This
// context can be used to collect attributes that will be added to log records.
// Check [AddAttrs] to add attributes to the context.
func WithAttrContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(attrContextKey{}).(*attrContext); ok {
		return ctx
	}
	return context.WithValue(ctx, attrContextKey{}, &attrContext{})
}

// AddAttrs adds attributes to the attribute context in the provided context.
// If the context does not have an attribute context, it will be ignored.
func AddAttrs(ctx context.Context, attrs ...slog.Attr) {
	if ac, ok := ctx.Value(attrContextKey{}).(*attrContext); ok {
		ac.attrs = append(ac.attrs, attrs...)
	}
}

// AddArgs adds key-value pairs as attributes to the attribute context in
// the provided context. Check [slog.Logger.Log] for more details on how to use.
func AddArgs(ctx context.Context, args ...any) {
	attrs := argsToAttrSlice(args)
	AddAttrs(ctx, attrs...)
}

// newAttrContextMiddleware creates a middleware that adds attributes from the
// attribute context to the log record.
func newAttrContextMiddleware() slogmulti.Middleware {
	return slogmulti.NewHandleInlineMiddleware(
		func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
			// If the context has an attribute context, add its attributes to the record.
			if ac, ok := ctx.Value(attrContextKey{}).(*attrContext); ok && len(ac.attrs) > 0 {
				record.AddAttrs(ac.attrs...)
			}

			return next(ctx, record)
		},
	)
}

// argsToAttrSlice converts a slice of any type to a slice of slog.Attr. This
// is copied from Go slog's internal implementation to ensure compatibility.
func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

// argsToAttr converts a slice of any type to a slog.Attr. This is copied from
// Go slog's internal implementation to ensure compatibility.
func argsToAttr(args []any) (slog.Attr, []any) {
	const badKey = "!BADKEY"

	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}
