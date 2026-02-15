package logger

import (
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

func newResource(args InitLoggerArgs) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(args.ServiceName),
			semconv.ServiceVersion(args.ServiceVersion),
			semconv.ServiceNamespace(args.ServiceNamespace),
			semconv.TelemetrySDKLanguageGo,
		),
	)
}
