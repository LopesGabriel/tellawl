package repositories

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/models"
)

type EmailNotificationTargetRepository interface {
	Upsert(ctx context.Context, target *models.EmailNotificationTarget) error
	List(ctx context.Context) ([]models.EmailNotificationTarget, error)
}
