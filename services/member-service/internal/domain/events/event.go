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
	MemberId  string    `json:"member_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

func (e MemberCreatedEvent) EventType() string {
	return "dev.lopesgabriel.member-service.member.created"
}
func (e MemberCreatedEvent) AggregateID() string   { return e.MemberId }
func (e MemberCreatedEvent) OccurredAt() time.Time { return e.Timestamp }
