package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/repositories"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgreSQLProcessedMessagesRepository struct {
	db        *sql.DB
	publisher events.EventPublisher
	tracer    trace.Tracer
}

func NewPostgreSQLProcessedMessagesRepository(db *sql.DB, publisher events.EventPublisher, tracer trace.Tracer) repositories.ProcessedMessagesRepository {
	if tracer == nil {
		tracer = tracing.GetTracer("github.com/lopesgabriel/tellawl/services/notifier/internal/infra/database")
	}

	return &postgreSQLProcessedMessagesRepository{
		db:        db,
		publisher: publisher,
		tracer:    tracer,
	}
}

func (r *postgreSQLProcessedMessagesRepository) Save(ctx context.Context, message *models.ProcessedMessage) error {
	ctx, span := r.tracer.Start(ctx, "SaveProcessedMessage", trace.WithAttributes(
		attribute.String("message.id", message.MessageID),
	))
	defer span.End()

	query := `
		INSERT INTO processed_messages (id, recipient, subject, sender, processed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, message.MessageID, message.Recepient, message.Subject, message.Sender, message.ProcessedAt)
	if err != nil {
		span.SetStatus(codes.Error, "failed to save processed message")
		span.RecordError(err)
		return err
	}

	if err := r.publisher.Publish(ctx, message.GetEvents()); err != nil {
		span.SetStatus(codes.Error, "failed to publish events")
		span.RecordError(err)
		return err
	}

	message.ClearEvents()
	span.SetStatus(codes.Ok, "success")
	return nil
}

func (r *postgreSQLProcessedMessagesRepository) Exists(ctx context.Context, messageID string) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "CheckProcessedMessageExists", trace.WithAttributes(
		attribute.String("message.id", messageID),
	))
	defer span.End()

	query := `
		SELECT processed_at
		FROM processed_messages
		WHERE id = $1
	`
	var processedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(&processedAt)
	if err != nil && err != sql.ErrNoRows {
		span.SetStatus(codes.Error, "failed to check if processed message exists")
		span.RecordError(err)
		return false, err
	}

	if processedAt.Valid {
		span.SetAttributes(attribute.String("processed_at", processedAt.Time.Format(time.RFC3339)))
	}

	span.SetStatus(codes.Ok, "success")
	return processedAt.Valid, nil
}
