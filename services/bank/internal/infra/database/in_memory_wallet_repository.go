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

type InMemoryWalletRepository struct {
	Items     []models.Wallet
	publisher ports.EventPublisher
}

func NewInMemoryWalletRepository(publisher ports.EventPublisher) *InMemoryWalletRepository {
	return &InMemoryWalletRepository{
		Items:     []models.Wallet{},
		publisher: publisher,
	}
}

func (r *InMemoryWalletRepository) Save(wallet *models.Wallet) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if wallet.Id == "" {
		wallet.Id = uuid.NewString()
	}

	// Handling update
	for i, w := range r.Items {
		if w.Id == wallet.Id {
			r.Items[i] = *wallet
			if err := r.publisher.Publish(ctx, wallet.Events()); err != nil {
				slog.Error("error publishing events", slog.String("error", err.Error()))
			}
			wallet.ClearEvents()
			return nil
		}
	}

	r.Items = append(r.Items, *wallet)
	if err := r.publisher.Publish(ctx, wallet.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}

	wallet.ClearEvents()
	return nil
}

func (r *InMemoryWalletRepository) FindById(id string) (*models.Wallet, error) {
	for _, wallet := range r.Items {
		if wallet.Id == id {
			return &wallet, nil
		}
	}

	return nil, repository.ErrWalletNotFound
}
