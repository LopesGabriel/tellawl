package repository

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
)

type MembersRepository interface {
	FindByID(ctx context.Context, id string) (*models.Member, error)
	FindByEmail(ctx context.Context, email string) (*models.Member, error)
	Upsert(ctx context.Context, member *models.Member) error
	Close() error
}

type Repositories struct {
	Members MembersRepository
}

func (r *Repositories) Close() error {
	if r.Members != nil {
		return r.Members.Close()
	}

	return nil
}
