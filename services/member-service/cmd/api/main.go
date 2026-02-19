package main

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/config"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/api"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/publisher"
	uc "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	configuration := config.InitAppConfigurations()

	// Telemetry initialization
	shutdown, err := initTelemetry(ctx, configuration)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	appLogger, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	// Init broker client
	var kafkaBroker broker.Broker
	if len(configuration.KafkaBrokers) > 0 {
		kafkaBroker, err = broker.NewKafkaBroker(broker.NewKafkaBrokerArgs{
			BootstrapServers: configuration.KafkaBrokers,
			Service:          configuration.ServiceName,
			Topic:            configuration.WalletTopic,
			Logger:           appLogger,
		})
		if err != nil {
			panic(err)
		}
		defer kafkaBroker.Close()
	}

	// Publisher initialization
	publisher := publisher.InitPublisher(kafkaBroker, appLogger)

	// Database initialization
	repos, err := database.InitDatabase(ctx, configuration, appLogger, publisher)
	if err != nil {
		panic(err)
	}
	defer repos.Close()

	// Use cases initialization
	usecases := uc.InitUseCases(uc.InitUseCasesArgs{
		JwtSecret: configuration.JwtSecret,
		Repos:     repos,
		Tracer:    tracing.GetTracer("github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"),
		Logger:    appLogger,
	})

	api := api.NewApiHandler(api.NewAPiArgs{
		Repositories: repos,
		Usecases:     usecases,
		Version:      configuration.Version,
		Tracer:       tracing.GetTracer("github.com/lopesgabriel/tellawl/services/member-service/internal/infra/api"),
		Logger:       appLogger,
	})
	err = api.Listen(ctx, configuration.Port)
	if err != nil {
		panic(err)
	}
}

func initTelemetry(ctx context.Context, appConfig *config.AppConfiguration) (func() error, error) {
	appLogger, err := logger.Init(ctx, logger.InitLoggerArgs{
		CollectorURL:     appConfig.OTELCollectorUrl,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.Version,
		Level:            slog.LevelDebug,
		LoggerProvider:   nil,
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
		if err := appLogger.Shutdown(ctx); err != nil {
			return err
		}
		if err := traceProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}
