package publisher

import (
	"context"

	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/config"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
)

func InitEventPublisher(ctx context.Context, config *config.AppConfiguration, appLogger *logger.AppLogger, broker broker.Broker) events.EventPublisher {
	if broker != nil {
		return NewKafkaPublisher(config, broker)
	}

	appLogger.Warn(ctx, "Nenhum broker de mensagens configurado. Usando InMemoryEventPublisher. Isso não é recomendado para produção.")
	return NewInMemoryEventPublisher(appLogger)
}
