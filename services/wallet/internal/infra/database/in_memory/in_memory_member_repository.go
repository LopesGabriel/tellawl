package inmemory

import (
	"context"
	"fmt"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/ports"
)

type InMemoryMemberRepository struct {
	Items     []models.Member
	publisher ports.EventPublisher
}

func NewInMemoryMemberRepository(publisher ports.EventPublisher) *InMemoryMemberRepository {
	return &InMemoryMemberRepository{
		Items:     []models.Member{},
		publisher: publisher,
	}
}

func (r InMemoryMemberRepository) FindByID(ctx context.Context, id string) (*models.Member, error) {
	var member models.Member

	for _, u := range r.Items {
		if u.Id == id {
			member = u
			break
		}
	}

	if member.Id == "" {
		return nil, fmt.Errorf("member not found")
	}

	return &member, nil
}

func (r InMemoryMemberRepository) FindByEmail(ctx context.Context, email string) (*models.Member, error) {
	var member models.Member

	for _, u := range r.Items {
		if u.Email == email {
			member = u
			break
		}
	}

	if member.Id == "" {
		return nil, fmt.Errorf("member not found")
	}

	return &member, nil
}

func (r *InMemoryMemberRepository) ValidateToken(ctx context.Context, token string) (*models.Member, error) {
	return nil, fmt.Errorf("token validation not implemented in in-memory repository")
}
