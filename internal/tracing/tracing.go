package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlphttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewTracingExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlphttp.New(ctx, otlphttp.WithEndpointURL("http://localhost:4318"))
}

func NewTracingProvider(exp sdktrace.SpanExporter, service string) *sdktrace.TracerProvider {

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
}
