package events

import (
	"time"
)

type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type WalletCreatedEvent struct {
	WalletId  string
	CreatorId string
	Name      string
	Timestamp time.Time
}

func (e WalletCreatedEvent) EventType() string     { return "com.tellawl.wallet.created" }
func (e WalletCreatedEvent) AggregateID() string   { return e.WalletId }
func (e WalletCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type WalletShared struct {
	WalletId  string
	UserId    string
	Timestamp time.Time
}

func (e WalletShared) EventType() string     { return "com.tellawl.wallet.shared" }
func (e WalletShared) AggregateID() string   { return e.WalletId }
func (e WalletShared) OccurredAt() time.Time { return e.Timestamp }

type TransactionRegisteredEvent struct {
	TransactionId string
	WalletId      string
	UserId        string
	Description   string
	Amount        map[string]int
	Type          string
	Timestamp     time.Time
}

func (e TransactionRegisteredEvent) EventType() string {
	return "com.tellawl.wallet.transaction.registered"
}
func (e TransactionRegisteredEvent) AggregateID() string   { return e.WalletId }
func (e TransactionRegisteredEvent) OccurredAt() time.Time { return e.Timestamp }
