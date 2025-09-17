package tracing

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func newExporter(ctx context.Context, args NewTraceProviderArgs) (*otlptrace.Exporter, error) {
	if args.CollectorURL == "" {
		return nil, fmt.Errorf("OTEL collector URL is required")
	}

	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(args.CollectorURL),
	}

	if !strings.Contains(args.CollectorURL, "https://") {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.New(ctx, options...)
}
