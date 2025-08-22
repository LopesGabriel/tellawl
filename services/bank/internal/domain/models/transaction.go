package models

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
)

type Transaction struct {
	Id string

	Amount    Monetary
	Category  Category
	CreatedBy User
	Type      TransactionType

	CreatedAt time.Time
	UpdatedAt *time.Time
}
