package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/config"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type kafkaEventPublisher struct {
	client broker.Broker
	tracer trace.Tracer
	logger *logger.AppLogger
}

func NewKafkaPublisher(appConfig *config.AppConfiguration) *kafkaEventPublisher {
	appLogger, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	broker, err := broker.NewKafkaBroker(broker.NewKafkaBrokerArgs{
		BootstrapServers: appConfig.KafkaBrokers,
		Service:          appConfig.ServiceName,
		Topic:            appConfig.KafkaTopic,
		Logger:           appLogger,
	})

	return &kafkaEventPublisher{
		client: broker,
		tracer: tracing.GetTracer("github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events/kafkaEventPublisher"),
		logger: appLogger,
	}
}

func (k *kafkaEventPublisher) Publish(ctx context.Context, events []events.DomainEvent) error {
	ctx, span := k.tracer.Start(ctx, "Publish")
	defer span.End()

	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			span.SetStatus(codes.Error, "failed to marshal event")
			span.RecordError(err)
			return err
		}

		eventId := uuid.NewString()

		err = k.client.Produce(ctx, &broker.Message{
			EventId:   eventId,
			EventType: event.EventType(),
			Key:       event.AggregateID(),
			Payload:   payload,
		})
		if err != nil {
			span.SetStatus(codes.Error, "failed to produce event")
			span.RecordError(err)
			return err
		}

		span.AddEvent(fmt.Sprintf("Produced event %s", eventId))
		k.logger.Info(ctx, "event published successfully",
			slog.Any("event.id", eventId),
			slog.String("event.type", event.EventType()),
			slog.String("event.aggregate_id", event.AggregateID()),
		)
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (k *kafkaEventPublisher) Close() error {
	return k.client.Close()
}
