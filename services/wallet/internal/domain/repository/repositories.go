package repository

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/ports"
	inmemory "github.com/lopesgabriel/tellawl/services/bank/internal/infra/database/in_memory"
)

type Repositories struct {
	User interface {
		FindByID(id string) (*models.User, error)
		FindByEmail(email string) (*models.User, error)
		Save(user *models.User) error
	}
	Wallet interface {
		FindById(id string) (*models.Wallet, error)
		FindByUserId(userId string) ([]models.Wallet, error)
		Save(wallet *models.Wallet) error
	}
}

func NewInMemory(publisher ports.EventPublisher) *Repositories {
	return &Repositories{
		User:   inmemory.NewInMemoryUserRepository(publisher),
		Wallet: inmemory.NewInMemoryWalletRepository(publisher),
	}
}
