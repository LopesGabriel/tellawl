package publisher

import (
	"context"
	"log/slog"
	"time"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/events"
)

type inMemoryEventPublisher struct {
	logger *logger.AppLogger
}

func InitInMemoryEventPublisher(logger *logger.AppLogger) events.EventPublisher {
	return &inMemoryEventPublisher{
		logger: logger,
	}
}

func (p *inMemoryEventPublisher) Publish(ctx context.Context, events []events.DomainEvent) error {
	for _, event := range events {
		p.logger.Info(
			ctx,
			"Evento publicado",
			slog.String("occurred_at", event.OccurredAt().UTC().Format(time.RFC3339)),
			slog.String("event_type", event.EventType()),
			slog.Any("event", event),
		)
	}

	return nil
}
