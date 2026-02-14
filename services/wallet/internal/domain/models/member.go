package models

import (
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
)

type Member struct {
	Id string

	FirstName string
	LastName  string
	Email     string

	CreatedAt time.Time
	UpdatedAt *time.Time

	events []events.DomainEvent
}

func (u *Member) AddEvent(event events.DomainEvent) {
	u.events = append(u.events, event)
}

func (u *Member) Events() []events.DomainEvent {
	return u.events
}

func (u *Member) ClearEvents() {
	u.events = nil
}
