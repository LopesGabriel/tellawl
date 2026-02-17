package publisher

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
)

type inMemoryEventPublisher struct {
	logger *logger.AppLogger
}

func NewInMemoryEventPublisher(logger *logger.AppLogger) *inMemoryEventPublisher {
	return &inMemoryEventPublisher{
		logger: logger,
	}
}

func (p inMemoryEventPublisher) Publish(ctx context.Context, events []events.DomainEvent) error {
	for _, event := range events {
		p.logger.Debug(
			ctx,
			"New Event published",
			slog.Any("occurred_at", event.OccurredAt()),
			slog.String("event_type", event.EventType()),
			slog.Any("event", event),
		)
	}

	return nil
}

func (p inMemoryEventPublisher) Close() error {
	return nil
}
