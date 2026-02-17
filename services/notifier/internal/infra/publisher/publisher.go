package publisher

import (
	"context"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/config"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/events"
)

func InitEventPublisher(ctx context.Context, config *config.AppConfiguration, appLogger *logger.AppLogger) events.EventPublisher {
	if len(config.KafkaBrokers) > 0 && config.KafkaTopic != "" {
		return NewKafkaPublisher(config)
	}

	appLogger.Warn(ctx, "Nenhum broker de mensagens configurado. Usando InMemoryEventPublisher. Isso não é recomendado para produção.")
	return InitInMemoryEventPublisher(appLogger)
}
