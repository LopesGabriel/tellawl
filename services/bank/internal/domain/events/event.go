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

func (e WalletCreatedEvent) EventType() string     { return "bank.wallet.created" }
func (e WalletCreatedEvent) AggregateID() string   { return e.WalletId }
func (e WalletCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type UserCreatedEvent struct {
	UserId    string
	FirstName string
	LastName  string
	Email     string
	Timestamp time.Time
}

func (e UserCreatedEvent) EventType() string     { return "bank.user.created" }
func (e UserCreatedEvent) AggregateID() string   { return e.UserId }
func (e UserCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type UserAddedToWalletEvent struct {
	WalletId  string
	UserId    string
	Timestamp time.Time
}

func (e UserAddedToWalletEvent) EventType() string     { return "bank.user.added_to_wallet" }
func (e UserAddedToWalletEvent) AggregateID() string   { return e.WalletId }
func (e UserAddedToWalletEvent) OccurredAt() time.Time { return e.Timestamp }

type TransactionRegisteredEvent struct {
	TransactionId string
	WalletId      string
	UserId        string
	Amount        map[string]int
	CategoryId    string
	Type          string
	Timestamp     time.Time
}

func (e TransactionRegisteredEvent) EventType() string     { return "bank.transaction.registered" }
func (e TransactionRegisteredEvent) AggregateID() string   { return e.WalletId }
func (e TransactionRegisteredEvent) OccurredAt() time.Time { return e.Timestamp }
