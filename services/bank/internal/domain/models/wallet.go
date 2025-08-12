package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/events"
)

type Wallet struct {
	ID        string
	CreatorID string
	Name      string
	Balance   Monetary
	Users     []User

	CreatedAt time.Time
	UpdatedAt *time.Time

	events []events.DomainEvent
}

func (w *Wallet) AddEvent(event events.DomainEvent) {
	w.events = append(w.events, event)
}

func (w *Wallet) Events() []events.DomainEvent {
	return w.events
}

func (w *Wallet) ClearEvents() {
	w.events = nil
}

func CreateNewWallet(name string, creator *User) *Wallet {
	wallet := &Wallet{
		ID:        uuid.NewString(),
		CreatorID: creator.ID,
		Name:      name,
		Balance: Monetary{
			Value:  0,
			Offset: 100,
		},
		Users:     []User{*creator},
		CreatedAt: time.Now(),
	}

	wallet.AddEvent(events.WalletCreatedEvent{
		WalletID:  wallet.ID,
		CreatorID: creator.ID,
		Name:      name,
		Timestamp: wallet.CreatedAt,
	})

	return wallet
}
