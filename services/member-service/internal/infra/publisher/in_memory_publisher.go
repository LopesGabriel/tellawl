package publisher

import (
	"context"
	"fmt"
	"time"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
)

type inMemoryEventPublisher struct {
}

func InitInMemoryEventPublisher() *inMemoryEventPublisher {
	return &inMemoryEventPublisher{}
}

func (p *inMemoryEventPublisher) Publish(ctx context.Context, events []events.DomainEvent) error {
	for _, event := range events {
		fmt.Printf("%s %s: %s\n", event.OccurredAt().UTC().Format(time.RFC3339), event.EventType(), event)
	}

	return nil
}
