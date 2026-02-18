package repositories

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/models"
)

type TelegramNotificationTargetRepository interface {
	Upsert(ctx context.Context, target *models.TelegramNotificationTarget) error
	List(ctx context.Context) ([]models.TelegramNotificationTarget, error)
}
