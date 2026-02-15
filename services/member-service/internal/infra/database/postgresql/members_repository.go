package postgresql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgreSQLMembersRepository struct {
	db        *sql.DB
	publisher events.EventPublisher
	tracer    trace.Tracer
	logger    *logger.AppLogger
}

func NewPostgreSQLMembersRepository(db *sql.DB, publisher events.EventPublisher) *postgreSQLMembersRepository {
	appLogger, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	return &postgreSQLMembersRepository{
		db:        db,
		publisher: publisher,
		tracer:    tracing.GetTracer("github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database/postgresql/postgreSQLMembersRepository"),
		logger:    appLogger,
	}
}

func (r *postgreSQLMembersRepository) FindByID(ctx context.Context, id string) (*models.Member, error) {
	ctx, span := r.tracer.Start(ctx, "FindByID", trace.WithAttributes(attribute.String("member.id", id)))
	defer span.End()

	query := `SELECT
		id, first_name, last_name, email, hashed_password, created_at, updated_at
	FROM members WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var member models.Member
	err := row.Scan(&member.Id, &member.FirstName, &member.LastName, &member.Email, &member.HashedPassword, &member.CreatedAt, &member.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound // Member not found
		}

		span.SetStatus(codes.Error, "failed to get member by id")
		span.RecordError(err)
		return nil, err
	}

	return &member, nil
}

func (r *postgreSQLMembersRepository) FindByEmail(ctx context.Context, email string) (*models.Member, error) {
	ctx, span := r.tracer.Start(ctx, "FindByEmail", trace.WithAttributes(attribute.String("member.email", email)))
	defer span.End()

	query := `SELECT
		id, first_name, last_name, email, hashed_password, created_at, updated_at
	FROM members WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var member models.Member
	err := row.Scan(&member.Id, &member.FirstName, &member.LastName, &member.Email, &member.HashedPassword, &member.CreatedAt, &member.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound // Member not found
		}

		span.SetStatus(codes.Error, "failed to get member by email")
		span.RecordError(err)
		return nil, err
	}

	return &member, nil
}

func (r *postgreSQLMembersRepository) Upsert(ctx context.Context, member *models.Member) error {
	ctx, span := r.tracer.Start(ctx, "Upsert", trace.WithAttributes(
		attribute.String("member.id", member.Id),
		attribute.String("member.email", member.Email),
	))
	defer span.End()

	if member.Id == "" {
		member.Id = uuid.NewString()
	}

	query := `INSERT INTO members (id, first_name, last_name, email, hashed_password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id)
	DO UPDATE SET
		first_name = EXCLUDED.first_name,
		last_name = EXCLUDED.last_name,
		email = EXCLUDED.email,
		hashed_password = EXCLUDED.hashed_password,
		updated_at = CURRENT_TIMESTAMP
	`

	span.AddEvent("Persisting in database")
	_, err := r.db.ExecContext(ctx, query, member.Id, member.FirstName, member.LastName, member.Email, member.HashedPassword, member.CreatedAt, member.UpdatedAt)
	if err != nil {
		span.SetStatus(codes.Error, "failed to upsert member")
		span.RecordError(err)
		return err
	}

	span.AddEvent("Publishing events")
	if err := r.publisher.Publish(ctx, member.Events()); err != nil {
		r.logger.Error(ctx, "error publishing events", slog.String("error", err.Error()))
	}
	member.ClearEvents()

	return nil
}
