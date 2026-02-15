package repository

import (
	"context"
	"database/sql"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	inmemory "github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database/in_memory"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database/postgresql"
)

type Repositories struct {
	Members interface {
		FindByID(ctx context.Context, id string) (*models.Member, error)
		FindByEmail(ctx context.Context, email string) (*models.Member, error)
		Upsert(ctx context.Context, member *models.Member) error
	}
}

func NewInMemory(publisher events.EventPublisher) *Repositories {
	return &Repositories{
		Members: inmemory.InitInMemoryMemberRepository(publisher),
	}
}

func NewPostgreSQL(db *sql.DB, publisher events.EventPublisher) *Repositories {
	return &Repositories{
		Members: postgresql.NewPostgreSQLMembersRepository(db, publisher),
	}
}
