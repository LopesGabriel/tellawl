package repositories

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/models"
)

type ProcessedMessagesRepository interface {
	Save(ctx context.Context, message *models.ProcessedMessage) error
	Exists(ctx context.Context, messageID string) (bool, error)
}
