package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/config"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func InitDatabase(ctx context.Context, appConfig *config.AppConfiguration, publisher events.EventPublisher) (*repository.Repositories, error) {
	var repo *repository.Repositories
	var httMemberRepo *HTTPMemberRepository

	appLogger, err := logger.GetLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if appConfig.MemberServiceUrl != "" {
		// Instrumented HTTP Client
		client := &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}

		repo, err := NewHTTPMemberRepository(appConfig.MemberServiceUrl, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP member repository: %w", err)
		}
		httMemberRepo = repo
	}

	if appConfig.DatabaseUrl != "" && appConfig.MigrationUrl != "" {
		db, err := NewPostgresClient(ctx, appConfig.DatabaseUrl)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}

		appLogger.Info(ctx, "Successfully connected to Postgres database")

		err = MigrateUp(appConfig.MigrationUrl, appConfig.DatabaseUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to apply database migration: %w", err)
		}

		repo = newPostgreSQL(db, publisher, httMemberRepo)
		return repo, nil
	}

	appLogger.Warn(ctx, "Using InMemory database. Data will not be persisted and will be lost on service restart.")

	return NewInMemory(publisher), nil
}

func NewInMemory(publisher events.EventPublisher) *repository.Repositories {
	return &repository.Repositories{
		Member: NewInMemoryMemberRepository(publisher),
		Wallet: NewInMemoryWalletRepository(publisher),
	}
}

func newPostgreSQL(db *sql.DB, publisher events.EventPublisher, memberRepo *HTTPMemberRepository) *repository.Repositories {
	return &repository.Repositories{
		Member: memberRepo,
		Wallet: NewPostgreSQLWalletRepository(db, publisher, memberRepo),
	}
}
