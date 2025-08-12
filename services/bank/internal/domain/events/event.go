package events

import "time"

type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type WalletCreatedEvent struct {
	WalletID  string
	CreatorID string
	Name      string
	Timestamp time.Time
}

func (e WalletCreatedEvent) EventType() string     { return "bank.wallet.created" }
func (e WalletCreatedEvent) AggregateID() string   { return e.WalletID }
func (e WalletCreatedEvent) OccurredAt() time.Time { return e.Timestamp }
