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

type MemberCreatedEvent struct {
	UserId    string
	FirstName string
	LastName  string
	Email     string
	Timestamp time.Time
}

func (e MemberCreatedEvent) EventType() string {
	return "dev.lopesgabriel.member-service.member.created"
}
func (e MemberCreatedEvent) AggregateID() string   { return e.UserId }
func (e MemberCreatedEvent) OccurredAt() time.Time { return e.Timestamp }
