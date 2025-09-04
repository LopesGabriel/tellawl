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

	Amount      Monetary
	CreatedBy   User
	Type        TransactionType
	Description string

	CreatedAt time.Time
	UpdatedAt *time.Time
}
