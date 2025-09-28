package main

import (
	"context"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/config"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/api"
	inmemoryevents "github.com/lopesgabriel/tellawl/services/member-service/internal/infra/events/in_memory"
	uc "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
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

	tracer := tracing.GetTracer(configuration.ServiceName)
	publisher := inmemoryevents.InitInMemoryEventPublisher()
	repos := repository.NewInMemory(publisher)
	usecases := uc.InitUseCases(uc.InitUseCasesArgs{
		JwtSecret: configuration.JwtSecret,
		Repos:     repos,
		Tracer:    tracer,
	})
	api := api.NewApiHandler(usecases, configuration.Version, tracer)
	err = api.Listen(ctx, configuration.Port)
	if err != nil {
		panic(err)
	}
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
