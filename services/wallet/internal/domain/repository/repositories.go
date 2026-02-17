package repository

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type Repositories struct {
	Member interface {
		FindByID(ctx context.Context, id string) (*models.Member, error)
		FindByEmail(ctx context.Context, email string) (*models.Member, error)
		ValidateToken(ctx context.Context, token string) (*models.Member, error)
	}
	Wallet interface {
		FindById(ctx context.Context, id string) (*models.Wallet, error)
		FindByUserId(ctx context.Context, userId string) ([]models.Wallet, error)
		Save(ctx context.Context, wallet *models.Wallet) error
	}
}
