package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type newTraceProviderArgs struct {
	Exporter *otlptrace.Exporter
	Version  string
}

func newTraceProvider(ctx context.Context, args newTraceProviderArgs) (*trace.TracerProvider, error) {
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("wallet"),
			semconv.ServiceVersion(args.Version),
			semconv.ServiceNamespace("tellawl"),
		),
	)

	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(args.Exporter),
		trace.WithResource(resource),
	)

	return traceProvider, nil
}
