package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
)

type inMemoryMemberRepository struct {
	items     map[string]models.Member
	publisher events.EventPublisher
}

func InitInMemoryMemberRepository(publisher events.EventPublisher) *inMemoryMemberRepository {
	return &inMemoryMemberRepository{
		items:     map[string]models.Member{},
		publisher: publisher,
	}
}

func (r inMemoryMemberRepository) FindByID(ctx context.Context, id string) (*models.Member, error) {
	var member models.Member

	member, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}

	return &member, nil
}

func (r inMemoryMemberRepository) FindByEmail(ctx context.Context, email string) (*models.Member, error) {
	var member models.Member

	for _, m := range r.items {
		if m.Email == email {
			member = m
			break
		}
	}

	if member.Id == "" {
		return nil, ErrNotFound
	}

	return &member, nil
}

func (r *inMemoryMemberRepository) Upsert(ctx context.Context, user *models.Member) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if user.Id == "" {
		user.Id = uuid.NewString()
	}

	existing, ok := r.items[user.Id]
	if ok {
		now := time.Now()
		user.CreatedAt = existing.CreatedAt
		user.UpdatedAt = &now
		r.items[user.Id] = *user
	} else {
		r.items[user.Id] = *user
	}

	if err := r.publisher.Publish(ctx, user.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	user.ClearEvents()

	return nil
}

func (r inMemoryMemberRepository) Close() error {
	return nil
}
