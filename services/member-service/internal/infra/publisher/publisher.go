package publisher

import (
	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
)

func InitPublisher(broker broker.Broker, appLogger *logger.AppLogger) events.EventPublisher {
	if broker == nil {
		return InitInMemoryEventPublisher()
	}

	return NewKafkaPublisher(broker, appLogger)
}
