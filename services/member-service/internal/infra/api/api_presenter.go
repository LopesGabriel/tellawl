package api

import (
	"time"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
)

type MemberAPIResponse struct {
	Id        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func toMemberAPIResponse(member *models.Member) MemberAPIResponse {
	return MemberAPIResponse{
		Id:        member.Id,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     member.Email,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func toMembersAPIResponse(members []models.Member) []MemberAPIResponse {
	response := make([]MemberAPIResponse, len(members))

	for i, member := range members {
		response[i] = toMemberAPIResponse(&member)
	}

	return response
}
