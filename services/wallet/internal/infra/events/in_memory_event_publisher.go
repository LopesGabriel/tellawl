package events

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
)

type InMemoryEventPublisher struct {
}

func (p InMemoryEventPublisher) Publish(ctx context.Context, events []events.DomainEvent) error {
	for _, event := range events {
		logger.Debug(
			ctx,
			"New Event published",
			slog.Any("occurred_at", event.OccurredAt()),
			slog.String("event_type", event.EventType()),
			slog.Any("event", event),
		)
	}

	return nil
}
