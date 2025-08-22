package presenter

import (
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type HTTPCategory struct {
	Id        string     `json:"id"`
	WalletId  string     `json:"wallet_id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewHTTPCategory(category models.Category) HTTPCategory {
	return HTTPCategory{
		Id:        category.Id,
		WalletId:  category.WalletId,
		Name:      category.Name,
		Type:      string(category.Type),
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (w HTTPCategory) ToJSON() []byte {
	data, err := json.Marshal(w)
	if err != nil {
		return []byte{}
	}

	return data
}
