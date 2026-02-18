package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/config"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/publisher"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	appConfig := config.InitAppConfigurations()

	// Telemetry initialization
	shutdown, err := initTelemetry(ctx, appConfig)
	if err != nil {
		fmt.Printf("failed to start telemetry: %v", err)
		panic(err)
	}
	defer shutdown()

	appLogger, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	var kafkaBroker broker.Broker
	if len(appConfig.KafkaBrokers) > 0 && appConfig.KafkaTopic != "" {
		kafkaBroker, err = broker.NewKafkaBroker(broker.NewKafkaBrokerArgs{
			BootstrapServers: appConfig.KafkaBrokers,
			Service:          appConfig.ServiceName,
			Topic:            appConfig.KafkaTopic,
			Logger:           appLogger,
		})
	}

	// Publisher initialization
	publisher := publisher.InitEventPublisher(ctx, appConfig, appLogger, kafkaBroker)
	defer publisher.Close()

	// Database initialization
	repos, err := database.InitDatabase(ctx, appConfig, publisher)
	if err != nil {
		appLogger.Fatal(ctx, "failed to initialize database", slog.String("error", err.Error()))
	}

	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		Repos:  repos,
		Logger: appLogger,
		Tracer: tracing.GetTracer("github.com/lopesgabriel/tellawl/service/wallet/internal/use-cases"),
	})
	apiHandler := controllers.NewAPIHandler(useCases, appConfig.Version)

	appLogger.Info(ctx, "Starting the API Server", slog.Int("port", appConfig.Port))
	apiHandler.Listen(appConfig.Port)
}

func initTelemetry(ctx context.Context, appConfig *config.AppConfiguration) (func() error, error) {
	appLogger, err := logger.Init(ctx, logger.InitLoggerArgs{
		CollectorURL:     appConfig.OTELCollectorUrl,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.Version,
		Level:            appConfig.LogLevel,
		LoggerProvider:   nil,
	})
	if err != nil {
		return nil, err
	}

	tracerProvider, err := tracing.Init(ctx, tracing.NewTraceProviderArgs{
		CollectorURL:     appConfig.OTELCollectorUrl,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.Version,
	})
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)

	return func() error {
		if err := appLogger.Shutdown(ctx); err != nil {
			return err
		}
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}
