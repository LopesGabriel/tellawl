package database

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// initialize a SQL client for postgresql
func NewPostgresClient(ctx context.Context, dbConnectionUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbConnectionUrl)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	return db, nil
}

func MigrateUp(migrationUrl, dbConnectionUrl string) error {
	m, err := migrate.New(migrationUrl, dbConnectionUrl)
	if err != nil {
		slog.Error("failed to load migration", slog.String("error", err.Error()))
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			slog.Error("failed to run migration", slog.String("error", err.Error()))
			return err
		}
	}

	return nil
}
