package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

type NewTraceProviderArgs struct {
	CollectorURL     string
	ServiceName      string
	ServiceNamespace string
	ServiceVersion   string
}

var Tracer trace.Tracer

func Init(ctx context.Context, args NewTraceProviderArgs) (*traceSdk.TracerProvider, error) {
	exporter, err := newExporter(ctx, args)
	if err != nil {
		return nil, err
	}

	traceProvider, err := newTraceProvider(exporter, args)
	if err != nil {
		return nil, err
	}

	return traceProvider, nil
}

func GetTracer(serviceName string) trace.Tracer {
	if Tracer == nil {
		Tracer = otel.Tracer(
			fmt.Sprintf("github.com/lopesgabriel/tellawl/service/%s", serviceName),
		)
	}

	return Tracer
}
