package presenter

import (
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type HTTPWallet struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	CreatorId    string            `json:"creator_id"`
	Balance      HTTPMonetary      `json:"balance"`
	Transactions []HTTPTransaction `json:"transactions"`
	Members      []HTTPMember      `json:"members"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    *time.Time        `json:"updated_at"`
}

func NewHTTPWallet(wallet models.Wallet) HTTPWallet {
	members := make([]HTTPMember, len(wallet.Members))
	transactions := make([]HTTPTransaction, len(wallet.Transactions))

	for i, member := range wallet.Members {
		members[i] = NewHTTPMember(member)
	}

	for i, transaction := range wallet.Transactions {
		transactions[i] = NewHTTPTransaction(transaction)
	}

	httpWallet := HTTPWallet{
		Id:           wallet.Id,
		Name:         wallet.Name,
		CreatorId:    wallet.CreatorId,
		Balance:      NewHTTPMonetary(wallet.Balance),
		Transactions: transactions,
		Members:      members,
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
