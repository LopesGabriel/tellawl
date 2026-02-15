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
	WalletId  string    `json:"wallet_id"`
	CreatorId string    `json:"creator_id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e WalletCreatedEvent) EventType() string     { return "com.tellawl.wallet.created" }
func (e WalletCreatedEvent) AggregateID() string   { return e.WalletId }
func (e WalletCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type WalletSharedEvent struct {
	WalletId  string    `json:"wallet_id"`
	MemberId  string    `json:"member_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (e WalletSharedEvent) EventType() string     { return "com.tellawl.wallet.shared" }
func (e WalletSharedEvent) AggregateID() string   { return e.WalletId }
func (e WalletSharedEvent) OccurredAt() time.Time { return e.Timestamp }

type TransactionRegisteredEvent struct {
	TransactionId string         `json:"transaction_id"`
	WalletId      string         `json:"wallet_id"`
	MemberId      string         `json:"member_id"`
	Description   string         `json:"description"`
	Amount        map[string]int `json:"amount"`
	Type          string         `json:"type"`
	Timestamp     time.Time      `json:"timestamp"`
}

func (e TransactionRegisteredEvent) EventType() string {
	return "com.tellawl.wallet.transaction.registered"
}
func (e TransactionRegisteredEvent) AggregateID() string   { return e.WalletId }
func (e TransactionRegisteredEvent) OccurredAt() time.Time { return e.Timestamp }
