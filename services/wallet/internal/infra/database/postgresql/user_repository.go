package postgresql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PostgreSQLUserRepository struct {
	db        *sql.DB
	publisher ports.EventPublisher
	tracer    trace.Tracer
}

func NewPostgreSQLUserRepository(db *sql.DB, publisher ports.EventPublisher) *PostgreSQLUserRepository {
	tracer := otel.Tracer("postgres-user-repository")

	return &PostgreSQLUserRepository{
		db:        db,
		publisher: publisher,
		tracer:    tracer,
	}
}

func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	ctx, span := r.tracer.Start(ctx, "PostgreSQLUserRepository.FindByID")
	defer span.End()

	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at 
			  FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	var updatedAt sql.NullTime

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	span.SetStatus(codes.Ok, "User found")
	return &user, nil
}

func (r *PostgreSQLUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := r.tracer.Start(ctx, "PostgreSQLUserRepository.FindByEmail")
	defer span.End()

	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at 
			  FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	var updatedAt sql.NullTime

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	span.SetStatus(codes.Ok, "User found")
	return &user, nil
}

func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *models.User) error {
	ctx, span := r.tracer.Start(ctx, "PostgreSQLUserRepository.Save")
	defer span.End()

	query := `INSERT INTO users (id, first_name, last_name, email, hashed_password, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  ON CONFLICT (id) DO UPDATE SET
			  first_name = EXCLUDED.first_name,
			  last_name = EXCLUDED.last_name,
			  email = EXCLUDED.email,
			  hashed_password = EXCLUDED.hashed_password,
			  updated_at = $8`

	now := time.Now()
	_, err := r.db.Exec(query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword,
		user.CreatedAt,
		user.UpdatedAt,
		now,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := r.publisher.Publish(ctx, user.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	user.ClearEvents()

	span.SetStatus(codes.Ok, "User saved")
	return err
}
