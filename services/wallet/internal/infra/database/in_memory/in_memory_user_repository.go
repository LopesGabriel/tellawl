package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/ports"
)

type InMemoryUserRepository struct {
	items     []models.User
	publisher ports.EventPublisher
}

func NewInMemoryUserRepository(publisher ports.EventPublisher) *InMemoryUserRepository {
	return &InMemoryUserRepository{
		items:     []models.User{},
		publisher: publisher,
	}
}

func (r InMemoryUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	for _, u := range r.items {
		if u.Id == id {
			user = u
			break
		}
	}

	if user.Id == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (r InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	for _, u := range r.items {
		if u.Email == email {
			user = u
			break
		}
	}

	if user.Id == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if user.Id == "" {
		user.Id = uuid.NewString()
	}

	if err := r.publisher.Publish(ctx, user.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	user.ClearEvents()

	r.items = append(r.items, *user)
	return nil
}
