package ports

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/events"
)

type EventPublisher interface {
	Publish(ctx context.Context, events []events.DomainEvent) error
}
