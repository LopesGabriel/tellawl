package database

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/ports"
)

type InMemoryUserRepository struct {
	Items     []models.User
	publisher ports.EventPublisher
}

func NewInMemoryUserRepository(publisher ports.EventPublisher) *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Items:     []models.User{},
		publisher: publisher,
	}
}

func (r *InMemoryUserRepository) FindByID(id string) (*models.User, error) {
	var user *models.User

	for _, u := range r.Items {
		if u.ID == id {
			user = &u
			break
		}
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
