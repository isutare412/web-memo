package trace

import (
	"context"
	"fmt"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerProvider atomic.Pointer[sdktrace.TracerProvider]

func Init(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}

	rsc, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
		))
	if err != nil {
		return fmt.Errorf("define resource: %w", err)
	}

	exp, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(cfg.OTLPGRPCEndpoint),
		otlptracegrpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("creating otlp grpc trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(rsc),
		sdktrace.WithRawSpanLimits(sdktrace.NewSpanLimits()),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SamplingRatio))),
	)

	tracerProvider.Store(tp)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

func Shutdown() error {
	tp := tracerProvider.Load()
	if tp == nil {
		return nil
	}

	if err := tp.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

func getTracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("github.com/isutare412/web-memo/api/internal/trace")
}
