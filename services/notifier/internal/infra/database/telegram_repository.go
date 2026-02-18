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

type postgresTelegramNotificationTargetRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func NewPostgreSQLTelegramNotificationTargetRepository(db *sql.DB, tracer trace.Tracer) repositories.TelegramNotificationTargetRepository {
	if tracer == nil {
		tracer = tracing.GetTracer("github.com/lopesgabriel/tellawl/services/notifier/internal/infra/database")
	}

	return &postgresTelegramNotificationTargetRepository{
		db:     db,
		tracer: tracer,
	}
}

func (r *postgresTelegramNotificationTargetRepository) Upsert(ctx context.Context, target *models.TelegramNotificationTarget) error {
	ctx, span := r.tracer.Start(ctx, "UpsertTelegramNotificationTarget", trace.WithAttributes(
		attribute.Int("telegram.chat_id", target.ChatID),
	))
	defer span.End()

	query := `
		INSERT INTO telegram_notification_targets (chat_id, nickname)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, target.ChatID, target.Nickname)
	if err != nil {
		span.SetStatus(codes.Error, "failed to save telegram notification target")
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (r *postgresTelegramNotificationTargetRepository) List(ctx context.Context) ([]models.TelegramNotificationTarget, error) {
	ctx, span := r.tracer.Start(ctx, "ListTelegramNotificationTargets", trace.WithAttributes())
	defer span.End()

	query := `
		SELECT chat_id, nickname
		FROM telegram_notification_targets
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, "failed to list telegram notification targets")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var targets []models.TelegramNotificationTarget
	for rows.Next() {
		var target models.TelegramNotificationTarget
		if err := rows.Scan(&target.ChatID, &target.Nickname); err != nil {
			span.SetStatus(codes.Error, "failed to scan telegram notification target")
			span.RecordError(err)
			return nil, err
		}
		targets = append(targets, target)
	}

	if err := rows.Err(); err != nil {
		span.SetStatus(codes.Error, "error iterating over telegram notification targets")
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return targets, nil
}
