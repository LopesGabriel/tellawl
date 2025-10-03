package presenter

import (
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type HTTPMember struct {
	Id        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewHTTPMember(member models.Member) HTTPMember {
	httpMember := HTTPMember{
		Id:        member.Id,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     member.Email,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}

	return httpMember
}

func (w HTTPMember) ToJSON() []byte {
	data, err := json.Marshal(w)
	if err != nil {
		return []byte{}
	}

	return data
}
