package presenter

import (
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type HTTPTransaction struct {
	Id              string       `json:"id"`
	Amount          HTTPMonetary `json:"amount"`
	CreatedBy       HTTPUser     `json:"created_by"`
	TransactionType string       `json:"transaction_type"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       *time.Time   `json:"updated_at"`
}

func NewHTTPTransaction(transaction models.Transaction) HTTPTransaction {
	httpTransaction := HTTPTransaction{
		Id:              transaction.Id,
		Amount:          NewHTTPMonetary(transaction.Amount),
		CreatedBy:       NewHTTPUser(transaction.CreatedBy),
		TransactionType: string(transaction.Type),
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	return httpTransaction
}

func (t HTTPTransaction) ToJSON() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		return []byte{}
	}

	return data
}
