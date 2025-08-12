package events

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/events"
)

type InMemoryEventPublisher struct {
}

func (p InMemoryEventPublisher) Publish(ctx context.Context, events []events.DomainEvent) error {
	for _, event := range events {
		slog.Debug("Event published", slog.Any("event", event))
	}

	return nil
}
