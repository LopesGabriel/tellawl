package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

type NewTraceProviderArgs struct {
	CollectorURL     string
	ServiceName      string
	ServiceNamespace string
	ServiceVersion   string
}

func Init(ctx context.Context, args NewTraceProviderArgs) (*traceSdk.TracerProvider, error) {
	exporter, err := newExporter(ctx, args)
	if err != nil {
		return nil, err
	}

	traceProvider, err := newTraceProvider(exporter, args)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})
	return traceProvider, nil
}

func GetTracer(serviceName string) trace.Tracer {
	if serviceName == "" {
		return otel.Tracer("github.com/lopesgabriel/tellawl")
	}

	return otel.Tracer(serviceName)
}
