package repository

import (
	"context"
	"database/sql"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/ports"
	inmemory "github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/in_memory"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/postgresql"
)

type Repositories struct {
	User interface {
		FindByID(ctx context.Context, id string) (*models.Member, error)
		FindByEmail(ctx context.Context, email string) (*models.Member, error)
		Save(ctx context.Context, member *models.Member) error
	}
	Wallet interface {
		FindById(ctx context.Context, id string) (*models.Wallet, error)
		FindByUserId(ctx context.Context, userId string) ([]models.Wallet, error)
		Save(ctx context.Context, wallet *models.Wallet) error
	}
}

func NewInMemory(publisher ports.EventPublisher) *Repositories {
	return &Repositories{
		User:   inmemory.NewInMemoryUserRepository(publisher),
		Wallet: inmemory.NewInMemoryWalletRepository(publisher),
	}
}

func NewPostgreSQL(db *sql.DB, publisher ports.EventPublisher) *Repositories {
	return &Repositories{
		User:   postgresql.NewPostgreSQLUserRepository(db, publisher),
		Wallet: postgresql.NewPostgreSQLWalletRepository(db, publisher),
	}
}
