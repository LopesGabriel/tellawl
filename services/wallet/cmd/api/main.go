package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/config"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	httpRepository "github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/http"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	// Database initialization
	db, err := database.NewPostgresClient(context.Background(), appConfig.DatabaseUrl)
	if err != nil {
		appLogger.Fatal(ctx, "failed to create the postgres client", slog.String("error", err.Error()))
	}
	err = db.Ping()
	if err != nil {
		appLogger.Fatal(ctx, "failed to ping database", slog.String("error", err.Error()))
	}
	appLogger.Info(ctx, "Database connected")

	err = database.MigrateUp(appConfig.MigrationUrl, appConfig.DatabaseUrl)
	if err != nil {
		appLogger.Fatal(ctx, "failed to apply database migration", slog.String("error", err.Error()))
	}

	// Publisher initialization
	publisher := events.NewKafkaPublisher(appConfig)
	defer publisher.Close()

	// Instrumented HTTP Client
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Repositories instance
	memberRepo, err := httpRepository.NewHTTPMemberRepository(appConfig.MemberServiceUrl, client)
	if err != nil {
		appLogger.Fatal(ctx, "failed to create HTTP member repository", slog.String("error", err.Error()))
	}
	repos := repository.NewPostgreSQL(db, publisher, memberRepo)

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
		Level:            slog.LevelDebug,
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
