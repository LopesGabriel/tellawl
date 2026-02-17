package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type InMemoryWalletRepository struct {
	items     []models.Wallet
	publisher events.EventPublisher
}

func NewInMemoryWalletRepository(publisher events.EventPublisher) *InMemoryWalletRepository {
	return &InMemoryWalletRepository{
		items:     []models.Wallet{},
		publisher: publisher,
	}
}

func (r *InMemoryWalletRepository) Save(ctx context.Context, wallet *models.Wallet) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if wallet.Id == "" {
		wallet.Id = uuid.NewString()
	}

	// Handling update
	for i, w := range r.items {
		if w.Id == wallet.Id {
			if err := r.publisher.Publish(ctx, wallet.Events()); err != nil {
				slog.Error("error publishing events", slog.String("error", err.Error()))
			}
			wallet.ClearEvents()

			r.items[i] = *wallet
			return nil
		}
	}

	if err := r.publisher.Publish(ctx, wallet.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	wallet.ClearEvents()

	r.items = append(r.items, *wallet)
	return nil
}

func (r InMemoryWalletRepository) FindById(ctx context.Context, id string) (*models.Wallet, error) {
	var wallet models.Wallet

	for _, w := range r.items {
		if w.Id == id {
			wallet = w
			break
		}
	}

	if wallet.Id == "" {
		return nil, fmt.Errorf("wallet not found")
	}

	return &wallet, nil
}

func (r InMemoryWalletRepository) FindByUserId(ctx context.Context, userId string) ([]models.Wallet, error) {
	var wallets []models.Wallet

	for _, wallet := range r.items {
		if wallet.IsMemberAllowedToRegisterTransactions(userId) {
			wallets = append(wallets, wallet)
		}
	}

	return wallets, nil
}
