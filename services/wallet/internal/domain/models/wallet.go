package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
)

type Wallet struct {
	Id        string
	CreatorId string
	Name      string
	Balance   Monetary

	Members      []Member
	Transactions []Transaction

	CreatedAt time.Time
	UpdatedAt *time.Time

	events []events.DomainEvent
}

func (w *Wallet) AddUser(member *Member) {
	w.Members = append(w.Members, *member)
	currentTime := time.Now()

	w.AddEvent(events.WalletShared{
		WalletId:  w.Id,
		UserId:    member.Id,
		Timestamp: currentTime,
	})
	w.UpdatedAt = &currentTime
}

func (w *Wallet) RegisterNewTransaction(amount Monetary, creator Member, transactionType TransactionType, description string) (*Transaction, error) {
	id := uuid.NewString()
	currentTime := time.Now()

	if amount.Value <= 0 {
		return nil, errors.New("invalid amount: must be greater than 0")
	}

	if !w.IsMemberAllowedToRegisterTransactions(creator.Id) {
		return nil, errors.New("user is not allowed to register transactions")
	}

	transaction := &Transaction{
		Id:          id,
		Amount:      amount,
		CreatedBy:   creator,
		Type:        transactionType,
		Description: description,
		CreatedAt:   currentTime,
	}

	w.Transactions = append(w.Transactions, *transaction)

	if transaction.Type == TransactionTypeDeposit {
		w.Balance = w.Balance.Sum(transaction.Amount)
	} else {
		w.Balance = w.Balance.Sub(transaction.Amount)
	}

	w.AddEvent(events.TransactionRegisteredEvent{
		TransactionId: transaction.Id,
		WalletId:      w.Id,
		UserId:        transaction.CreatedBy.Id,
		Amount: map[string]int{
			"value":  transaction.Amount.Value,
			"offset": transaction.Amount.Offset,
		},
		Type:        string(transaction.Type),
		Timestamp:   currentTime,
		Description: description,
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

func CreateNewWallet(name string, creator *Member) *Wallet {
	walletId := uuid.NewString()
	wallet := &Wallet{
		Id:        walletId,
		CreatorId: creator.Id,
		Name:      name,
		Balance: Monetary{
			Value:  0,
			Offset: 100,
		},
		Members:      []Member{*creator},
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

func (w *Wallet) IsMemberAllowedToRegisterTransactions(memberId string) bool {
	if w.CreatorId == memberId {
		return true
	}

	for _, member := range w.Members {
		if member.Id == memberId {
			return true
		}
	}

	return false
}
