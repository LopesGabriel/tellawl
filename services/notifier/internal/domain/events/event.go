package events

import (
	"context"
	"time"
)

type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type EventPublisher interface {
	Publish(ctx context.Context, events []DomainEvent) error
}

type EmailReceivedEvent struct {
	MessageId string    `json:"message_id"`
	Recipient string    `json:"recipient"`
	Sender    string    `json:"sender"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

func (e EmailReceivedEvent) EventType() string {
	return "dev.lopesgabriel.notifier.email.received"
}
func (e EmailReceivedEvent) AggregateID() string   { return e.MessageId }
func (e EmailReceivedEvent) OccurredAt() time.Time { return e.Timestamp }
