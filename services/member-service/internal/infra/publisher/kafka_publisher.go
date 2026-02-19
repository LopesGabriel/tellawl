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
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type kafkaEventPublisher struct {
	client broker.Broker
	tracer trace.Tracer
	logger *logger.AppLogger
}

func NewKafkaPublisher(broker broker.Broker, appLogger *logger.AppLogger) *kafkaEventPublisher {
	return &kafkaEventPublisher{
		client: broker,
		tracer: tracing.GetTracer("github.com/lopesgabriel/tellawl/services/member-service/internal/infra/publisher/kafkaEventPublisher"),
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
		k.logger.Info(ctx, "event published successfully", slog.Any("id", eventId), slog.String("type", event.EventType()))
	}

	return nil
}
