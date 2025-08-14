package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/events"
)

type Wallet struct {
	Id        string
	CreatorId string
	Name      string
	Balance   Monetary
	Users     []User

	CreatedAt time.Time
	UpdatedAt *time.Time

	events []events.DomainEvent
}

func (w *Wallet) AddUser(user *User) {
	w.Users = append(w.Users, *user)
	currentTime := time.Now()

	w.AddEvent(events.UserAddedToWalletEvent{
		WalletId:  w.Id,
		UserId:    user.Id,
		Timestamp: currentTime,
	})
	w.UpdatedAt = &currentTime
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
		Id:        uuid.NewString(),
		CreatorId: creator.Id,
		Name:      name,
		Balance: Monetary{
			Value:  0,
			Offset: 100,
		},
		Users:     []User{*creator},
		CreatedAt: time.Now(),
	}

	wallet.AddEvent(events.WalletCreatedEvent{
		WalletId:  wallet.Id,
		CreatorId: creator.Id,
		Name:      name,
		Timestamp: wallet.CreatedAt,
	})

	return wallet
}
