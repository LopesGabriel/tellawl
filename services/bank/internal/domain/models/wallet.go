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
	Categories   []Category

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

func (w *Wallet) RegisterNewTransaction(amount Monetary, creator User, transactionType TransactionType, categoryId string, description string) (*Transaction, error) {
	id := uuid.NewString()
	currentTime := time.Now()

	if amount.Value <= 0 {
		return nil, errors.New("invalid amount: must be greater than 0")
	}

	if !w.IsUserAllowedToRegisterTransactions(creator.Id) {
		return nil, errors.New("user is not allowed to register transactions")
	}

	var category Category
	for _, c := range w.Categories {
		if c.Id == categoryId {
			category = c
			break
		}
	}

	if category.Id == "" {
		return nil, errors.New("category not found")
	}

	transaction := &Transaction{
		Id:          id,
		Amount:      amount,
		Category:    category,
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
		CategoryId:  category.Id,
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

func CreateNewWallet(name string, creator *User) *Wallet {
	walletId := uuid.NewString()
	wallet := &Wallet{
		Id:        walletId,
		CreatorId: creator.Id,
		Name:      name,
		Balance: Monetary{
			Value:  0,
			Offset: 100,
		},
		Users:        []User{*creator},
		Transactions: []Transaction{},
		Categories:   createDefaultCategories(walletId),
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

func (w *Wallet) AddCustomCategory(name string) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category := CreateCustomCategory(w.Id, name)
	w.Categories = append(w.Categories, *category)
	currentTime := time.Now()
	w.UpdatedAt = &currentTime

	return category, nil
}

func (w *Wallet) UpdateCategory(categoryId, name string) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	for i, category := range w.Categories {
		if category.Id == categoryId {
			if category.Type == CategoryTypeDefault {
				return nil, errors.New("cannot update default categories")
			}
			w.Categories[i].Name = name
			currentTime := time.Now()
			w.Categories[i].UpdatedAt = &currentTime
			w.UpdatedAt = &currentTime
			return &w.Categories[i], nil
		}
	}

	return nil, errors.New("category not found")
}

func (w *Wallet) DeleteCategory(categoryId string) error {
	for i, category := range w.Categories {
		if category.Id == categoryId {
			if category.Type == CategoryTypeDefault {
				return errors.New("cannot delete default categories")
			}
			w.Categories = append(w.Categories[:i], w.Categories[i+1:]...)
			currentTime := time.Now()
			w.UpdatedAt = &currentTime
			return nil
		}
	}

	return errors.New("category not found")
}

func (w *Wallet) GetCategories() []Category {
	return w.Categories
}

func (w *Wallet) IsUserAllowedToRegisterTransactions(userId string) bool {
	if w.CreatorId == userId {
		return true
	}

	for _, user := range w.Users {
		if user.Id == userId {
			return true
		}
	}

	return false
}

func createDefaultCategories(walletId string) []Category {
	defaultNames := GetDefaultCategories()
	categories := make([]Category, len(defaultNames))

	for i, name := range defaultNames {
		categories[i] = *CreateDefaultCategory(walletId, name)
	}

	return categories
}
