package presenter

import (
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type HTTPWallet struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	CreatorId    string            `json:"creator_id"`
	Balance      HTTPMonetary      `json:"balance"`
	Transactions []HTTPTransaction `json:"transactions"`
	Users        []HTTPUser        `json:"users"`
	Categories   []HTTPCategory    `json:"categories"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    *time.Time        `json:"updated_at"`
}

func NewHTTPWallet(wallet models.Wallet) HTTPWallet {
	users := make([]HTTPUser, len(wallet.Users))
	transactions := make([]HTTPTransaction, len(wallet.Transactions))
	categories := make([]HTTPCategory, len(wallet.Categories))

	for i, user := range wallet.Users {
		users[i] = NewHTTPUser(user)
	}

	for i, transaction := range wallet.Transactions {
		transactions[i] = NewHTTPTransaction(transaction)
	}

	for i, category := range wallet.Categories {
		categories[i] = NewHTTPCategory(category)
	}

	httpWallet := HTTPWallet{
		Id:           wallet.Id,
		Name:         wallet.Name,
		CreatorId:    wallet.CreatorId,
		Balance:      NewHTTPMonetary(wallet.Balance),
		Transactions: transactions,
		Categories:   categories,
		Users:        users,
		CreatedAt:    wallet.CreatedAt,
		UpdatedAt:    wallet.UpdatedAt,
	}

	return httpWallet
}

func (w HTTPWallet) ToJSON() []byte {
	data, err := json.Marshal(w)
	if err != nil {
		return []byte{}
	}

	return data
}
