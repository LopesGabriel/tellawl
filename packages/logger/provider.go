package logger

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func newProvider(ctx context.Context, resource *resource.Resource, args InitLoggerArgs) (*log.LoggerProvider, error) {
	options := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(args.CollectorURL),
	}

	if !strings.Contains(args.CollectorURL, "grpcs://") {
		options = append(options, otlploggrpc.WithInsecure())
	}

	exporter, err := otlploggrpc.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	processor := log.NewBatchProcessor(exporter)
	provider := log.NewLoggerProvider(
		log.WithResource(resource),
		log.WithProcessor(processor),
	)

	return provider, nil
}
