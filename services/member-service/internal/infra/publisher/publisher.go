package publisher

import (
	"context"

	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
)

func InitPublisher(ctx context.Context, broker broker.Broker, appLogger *logger.AppLogger) events.EventPublisher {
	if broker == nil {
		appLogger.Warn(ctx, "No broker configured, using in-memory event publisher. Events will not be persisted or shared across instances.")
		return InitInMemoryEventPublisher()
	}

	return NewKafkaPublisher(broker, appLogger)
}
