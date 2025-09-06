package presenter

import (
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type HTTPUser struct {
	Id        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewHTTPUser(user models.User) HTTPUser {
	httpUser := HTTPUser{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return httpUser
}

func (w HTTPUser) ToJSON() []byte {
	data, err := json.Marshal(w)
	if err != nil {
		return []byte{}
	}

	return data
}
