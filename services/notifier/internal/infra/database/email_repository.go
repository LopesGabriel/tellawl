package database

import (
	"context"
	"database/sql"

	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/repositories"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgresEmailNotificationTargetRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func NewPostgreSQLEmailNotificationTargetRepository(db *sql.DB, tracer trace.Tracer) repositories.EmailNotificationTargetRepository {
	if tracer == nil {
		tracer = tracing.GetTracer("github.com/lopesgabriel/tellawl/services/notifier/internal/infra/database")
	}

	return &postgresEmailNotificationTargetRepository{
		db:     db,
		tracer: tracer,
	}
}

func (r *postgresEmailNotificationTargetRepository) Upsert(ctx context.Context, target *models.EmailNotificationTarget) error {
	ctx, span := r.tracer.Start(ctx, "UpsertEmailNotificationTarget", trace.WithAttributes(
		attribute.String("email.address", target.Email),
	))
	defer span.End()

	query := `
		INSERT INTO email_notification_targets (email, name)
		VALUES ($1, $2)
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
	`
	_, err := r.db.ExecContext(ctx, query, target.Email, target.Name)
	if err != nil {
		span.SetStatus(codes.Error, "failed to upsert email notification target")
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (r *postgresEmailNotificationTargetRepository) List(ctx context.Context) ([]models.EmailNotificationTarget, error) {
	ctx, span := r.tracer.Start(ctx, "ListEmailNotificationTargets")
	defer span.End()

	query := `
		SELECT id, email, name
		FROM email_notification_targets
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, "failed to list email notification targets")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var targets []models.EmailNotificationTarget
	for rows.Next() {
		var target models.EmailNotificationTarget
		if err := rows.Scan(&target.ID, &target.Email, &target.Name); err != nil {
			span.SetStatus(codes.Error, "failed to scan email notification target")
			span.RecordError(err)
			return nil, err
		}
		targets = append(targets, target)
	}

	if err := rows.Err(); err != nil {
		span.SetStatus(codes.Error, "error iterating over email notification targets")
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return targets, nil
}
