package tracing

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

func newTraceProvider(exporter *otlptrace.Exporter, args NewTraceProviderArgs) (*trace.TracerProvider, error) {
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(args.ServiceName),
		semconv.ServiceVersion(args.ServiceVersion),
		semconv.ServiceNamespace(args.ServiceNamespace),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKNameKey.String("opentelemetry"),
	)

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	return traceProvider, nil
}
