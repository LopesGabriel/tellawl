package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/events"
)

type Wallet struct {
	Id        string
	CreatorId string
	Name      string
	Balance   Monetary

	Users        []User
	Transactions []Transaction

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

func (w *Wallet) RegisterNewTransaction(amount Monetary, creator User, transactionType TransactionType) (*Transaction, error) {
	id := uuid.NewString()
	currentTime := time.Now()

	if amount.Value <= 0 {
		return nil, errors.New("invalid amount: must be greater than 0")
	}

	if !w.IsUserAllowedToRegisterTransactions(creator.Id) {
		return nil, errors.New("user is not allowed to register transactions")
	}

	transaction := &Transaction{
		Id:        id,
		Amount:    amount,
		CreatedBy: creator,
		Type:      transactionType,
		CreatedAt: currentTime,
	}

	w.Transactions = append(w.Transactions, *transaction)

	w.Balance = w.Balance.Sum(transaction.Amount)

	w.AddEvent(events.TransactionRegisteredEvent{
		TransactionId: transaction.Id,
		WalletId:      w.Id,
		UserId:        transaction.CreatedBy.Id,
		Amount: map[string]int{
			"value":  transaction.Amount.Value,
			"offset": transaction.Amount.Offset,
		},
		Type:      string(transaction.Type),
		Timestamp: currentTime,
	})

	return transaction, nil
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
		Users:        []User{*creator},
		Transactions: []Transaction{},
		CreatedAt:    time.Now(),
	}

	wallet.AddEvent(events.WalletCreatedEvent{
		WalletId:  wallet.Id,
		CreatorId: creator.Id,
		Name:      name,
		Timestamp: wallet.CreatedAt,
	})

	return wallet
}

func (w *Wallet) IsUserAllowedToRegisterTransactions(userId string) bool {
	for _, user := range w.Users {
		if user.Id == userId {
			return true
		}
	}

	return false
}
