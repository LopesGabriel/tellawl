package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/config"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/api"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/events/kafka"
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

	// Database initialization
	db, err := initDatabase(ctx, configuration, appLogger)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Broker initialization
	publisher := kafka.NewKafkaPublisher(configuration)

	// Repositories initialization
	repos := repository.NewPostgreSQL(db, publisher)

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

func initDatabase(ctx context.Context, appConfig *config.AppConfiguration, appLogger *logger.AppLogger) (*sql.DB, error) {
	// Initialize the PostgreSQL client
	db, err := database.NewPostgresClient(context.Background(), appConfig.DatabaseUrl)
	if err != nil {
		appLogger.Fatal(ctx, "failed to create the postgres client", slog.String("error", err.Error()))
	}

	// Ping the database to ensure the connection is established
	err = db.Ping()
	if err != nil {
		appLogger.Fatal(ctx, "failed to ping database", slog.String("error", err.Error()))
	}
	appLogger.Info(ctx, "Connected to database")

	// Apply database migrations
	err = database.MigrateUp(appConfig.MigrationUrl, appConfig.DatabaseUrl)
	if err != nil {
		appLogger.Fatal(ctx, "failed to apply database migration", slog.String("error", err.Error()))
	}

	return db, nil
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
