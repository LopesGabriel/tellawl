package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/config"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	appConfig := config.InitAppConfigurations()

	shutdown, err := initTelemetry(ctx, appConfig)
	if err != nil {
		fmt.Printf("failed to start telemetry: %v", err)
		panic(err)
	}
	defer shutdown()
	logger.Info(ctx, "Telemetry started")

	db, err := database.NewPostgresClient(context.Background(), appConfig.DatabaseUrl)
	if err != nil {
		logger.Fatal(ctx, "failed to create the postgres client", slog.String("error", err.Error()))
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal(ctx, "failed to ping database", slog.String("error", err.Error()))
	}
	logger.Info(ctx, "Database connected")

	err = database.MigrateUp(appConfig.MigrationUrl, appConfig.DatabaseUrl)
	if err != nil {
		logger.Fatal(ctx, "failed to apply database migration", slog.String("error", err.Error()))
	}

	publisher := events.InMemoryEventPublisher{}

	repos := repository.NewPostgreSQL(db, publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: os.Getenv("JWT_SECRET"),
		Repos:     repos,
	})
	apiHandler := controllers.NewAPIHandler(useCases, appConfig.Version)

	logger.Info(ctx, "Starting the API Server", slog.Int("port", appConfig.Port))
	apiHandler.Listen(appConfig.Port)
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
		if err := logProvider.Shutdown(ctx); err != nil {
			return err
		}
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}
