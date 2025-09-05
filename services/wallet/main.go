package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

const VERSION = "1.0.0"
const JWT_SECRET = "ex4mpl3-Secr3t"
const PORT = 8080

func main() {
	slog.Info("Starting Wallet Service", slog.String("version", VERSION))
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	dbConnectionString := os.Getenv("POSTGRESQL_URL")

	slog.Info("Starting database interface")
	db, err := database.NewPostgresClient(context.Background(), dbConnectionString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	slog.Info("Successfully connected to database")

	err = database.MigrateUp(os.Getenv("MIGRATE_URL"), dbConnectionString)
	if err != nil {
		panic(err)
	}

	publisher := events.InMemoryEventPublisher{}

	repos := repository.NewPostgreSQL(db, publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: JWT_SECRET,
		Repos:     repos,
	})
	apiHandler := controllers.NewAPIHandler(useCases, VERSION)

	slog.Info(fmt.Sprintf("API Server listenig on port %d", PORT))
	apiHandler.Listen(PORT)
}
