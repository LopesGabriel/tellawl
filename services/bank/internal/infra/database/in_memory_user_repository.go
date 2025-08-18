package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/ports"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
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
		if u.Id == id {
			user = &u
			break
		}
	}

	if user == nil {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*models.User, error) {
	var user *models.User

	for _, u := range r.Items {
		if u.Email == email {
			user = &u
			break
		}
	}

	if user == nil {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) Save(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if user.Id == "" {
		user.Id = uuid.NewString()
	}

	if err := r.publisher.Publish(ctx, user.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	user.ClearEvents()

	r.Items = append(r.Items, *user)
	return nil
}
