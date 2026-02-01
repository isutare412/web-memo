package trace

import (
	"context"
	"net/http"
	"path"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// AutoSpan generates a span with caller function name. If you need to use
// custom function name, use [StartSpan] instead.
func AutoSpan(ctx context.Context, opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return createSpan(ctx, funcName(2), opts...)
}

// StartSpan generates a span with given name.
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return createSpan(ctx, name, opts...)
}

func createSpan(ctx context.Context, name string, opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return getTracer().Start(ctx, name, opts...)
}

// AutoError record error with caller function name. If you need to use
// custom function name, use [RecordError] instead.
func AutoError(span trace.Span, err error, opts ...trace.EventOption) {
	recordError(span, err, funcName(2), opts...)
}

// RecordError record error to span with given name.
func RecordError(span trace.Span, err error, name string, opts ...trace.EventOption) {
	recordError(span, err, name, opts...)
}

func recordError(span trace.Span, err error, name string, opts ...trace.EventOption) {
	span.SetStatus(codes.Error, name)
	span.RecordError(err, opts...)
}

// SpanFromContext returns the current Span from ctx.
//
// If no Span is currently set in ctx an implementation of a Span that
// performs no operations is returned.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func CallerName() string {
	return funcName(2)
}

func funcName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "UnknownFunc"
	}
	caller := runtime.FuncForPC(pc)
	if caller == nil {
		return "UnknownFunc"
	}
	return path.Base(caller.Name())
}

func InjectToHeader(ctx context.Context, h http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(h))
}

func InjectToMap(ctx context.Context, m map[string]string) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(m))
}

func ExtractFromHTTPHeader(ctx context.Context, h http.Header) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(h))
}

func ExtractFromMap(ctx context.Context, m map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m))
}
