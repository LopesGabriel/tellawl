package main

import (
	"context"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/config"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	configuration := config.InitAppConfigurations()
	shutdown, err := initTelemetry(ctx, configuration)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	<-ctx.Done()
}

func initTelemetry(ctx context.Context, appConfig *config.AppConfiguration) (func() error, error) {
	logProvider, err := logger.Init(ctx, logger.InitLoggerArgs{
		CollectorURL:     appConfig.OTELCollectorUrl,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.Version,
	})
	if err != nil {
		return nil, err
	}

	traceProvider, err := tracing.Init(ctx, tracing.NewTraceProviderArgs{
		CollectorURL:     appConfig.OTELCollectorUrl,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.Version,
	})
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(traceProvider)

	return func() error {
		if err := logProvider.Shutdown(ctx); err != nil {
			return err
		}
		if err := traceProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}
