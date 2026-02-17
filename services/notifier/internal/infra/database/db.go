package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/exaring/otelpgx"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"go.opentelemetry.io/otel"
)

// initialize a SQL client for postgresql
func NewPostgresClient(ctx context.Context, dbConnectionUrl string) (*sql.DB, error) {
	appLogger, err := logger.GetLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to get logger: %w", err)
	}

	cfg, err := pgxpool.ParseConfig(dbConnectionUrl)
	if err != nil {
		err = fmt.Errorf("create connection pool: %w", err)
		appLogger.Error(ctx, "failed to create connection pool", slog.String("error", err.Error()))
		return nil, err
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTracerProvider(otel.GetTracerProvider()))

	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	db := stdlib.OpenDBFromPool(pool)
	if err != nil {
		appLogger.Error(ctx, "failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	return db, nil
}

func MigrateUp(migrationUrl, dbConnectionUrl string) error {
	m, err := migrate.New(migrationUrl, dbConnectionUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}
