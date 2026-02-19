package database

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/config"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
)

func InitDatabase(ctx context.Context, appConfig *config.AppConfiguration, appLogger *logger.AppLogger, publisher events.EventPublisher) (*repository.Repositories, error) {
	repositories := &repository.Repositories{}

	if appConfig.DatabaseUrl == "" || appConfig.MigrationUrl == "" {
		repositories.Members = InitInMemoryMemberRepository(publisher)
		appLogger.Warn(ctx, "Using InMemory Database! Not recommended for production environment")
		return repositories, nil
	}

	// Initialize the PostgreSQL client
	db, err := NewPostgresClient(context.Background(), appConfig.DatabaseUrl)
	if err != nil {
		appLogger.Fatal(ctx, "failed to create the postgres client", slog.String("error", err.Error()))
	}

	// Ping the database to ensure the connection is established
	err = db.Ping()
	if err != nil {
		appLogger.Fatal(ctx, "failed to ping database", slog.String("error", err.Error()))
	}

	appLogger.Info(ctx, "Connected to postgresql database")

	// Apply database migrations
	err = migrateUp(appConfig.MigrationUrl, appConfig.DatabaseUrl)
	if err != nil {
		appLogger.Fatal(ctx, "failed to apply database migration", slog.String("error", err.Error()))
	}

	repositories.Members = NewPostgreSQLMembersRepository(db, publisher, appLogger)

	return repositories, nil
}
