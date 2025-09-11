package telemetry

import (
	"context"
	"log/slog"
	"strings"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/core"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func InitTelemetry(ctx context.Context, configuration *core.Configuration) func() {
	exp, err := newTraceExporter(ctx, configuration)
	if err != nil {
		slog.Error("failed to create new trace exporter", "error", err)
		panic(err)
	}

	traceProvider, err := newTraceProvider(ctx, newTraceProviderArgs{Exporter: exp, Version: configuration.Version})
	if err != nil {
		slog.Error("failed to create new trace provider", "error", err)
		panic(err)
	}
	otel.SetTracerProvider(traceProvider)

	return func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown trace provider", "error", err)
		}
	}
}

func newTraceExporter(ctx context.Context, configuration *core.Configuration) (*otlptrace.Exporter, error) {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpointURL(configuration.OTELCollectorUrl),
	}

	if !strings.Contains(configuration.OTELCollectorUrl, "https://") {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.New(ctx, options...)
}
