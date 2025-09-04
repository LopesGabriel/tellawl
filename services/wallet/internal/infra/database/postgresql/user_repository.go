package postgresql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/ports"
)

type PostgreSQLUserRepository struct {
	db        *sql.DB
	publisher ports.EventPublisher
}

func NewPostgreSQLUserRepository(db *sql.DB, publisher ports.EventPublisher) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{
		db:        db,
		publisher: publisher,
	}
}

func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at 
			  FROM users WHERE id = $1`

	row := r.db.QueryRow(query, id)

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
		return nil, err
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	return &user, nil
}

func (r *PostgreSQLUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, email, hashed_password, created_at, updated_at 
			  FROM users WHERE email = $1`

	row := r.db.QueryRow(query, email)

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
		return nil, err
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	return &user, nil
}

func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *models.User) error {
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
		return err
	}

	if err := r.publisher.Publish(ctx, user.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	user.ClearEvents()

	return err
}
