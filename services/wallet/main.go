package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/core"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("Starting Wallet Service")

	appConfig := core.InitAppConfigurations()

	slog.Info("Starting database interface")
	db, err := database.NewPostgresClient(context.Background(), appConfig.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	slog.Info("Successfully connected to database")

	err = database.MigrateUp(appConfig.MigrationUrl, appConfig.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	publisher := events.InMemoryEventPublisher{}

	repos := repository.NewPostgreSQL(db, publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: os.Getenv("JWT_SECRET"),
		Repos:     repos,
	})
	apiHandler := controllers.NewAPIHandler(useCases, appConfig.Version)

	slog.Info(
		fmt.Sprintf("API Server listenig on port %d", appConfig.Port),
		slog.String("version", appConfig.Version),
	)
	apiHandler.Listen(appConfig.Port)
}
