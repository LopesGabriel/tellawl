package inmemory

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
)

type inMemoryMemberRepository struct {
	items     []models.Member
	publisher events.EventPublisher
}

func InitInMemoryMemberRepository(publisher events.EventPublisher) *inMemoryMemberRepository {
	return &inMemoryMemberRepository{
		items:     []models.Member{},
		publisher: publisher,
	}
}

func (r inMemoryMemberRepository) FindByID(ctx context.Context, id string) (*models.Member, error) {
	var member models.Member

	for _, u := range r.items {
		if u.Id == id {
			member = u
			break
		}
	}

	if member.Id == "" {
		return nil, database.ErrNotFound
	}

	return &member, nil
}

func (r inMemoryMemberRepository) FindByEmail(ctx context.Context, email string) (*models.Member, error) {
	var member models.Member

	for _, u := range r.items {
		if u.Email == email {
			member = u
			break
		}
	}

	if member.Id == "" {
		return nil, database.ErrNotFound
	}

	return &member, nil
}

func (r *inMemoryMemberRepository) Save(ctx context.Context, user *models.Member) error {
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
