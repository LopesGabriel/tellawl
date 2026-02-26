package listener

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/repositories"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/infra/email"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type kafkaListener struct {
	broker                      broker.Broker
	processedMessagesRepository repositories.ProcessedMessagesRepository
	emailRepo                   repositories.EmailNotificationTargetRepository
	logger                      *logger.AppLogger
	tracer                      trace.Tracer
	topic                       string
	emailClient                 *email.Client
}

type NewKafkaListenerParams struct {
	Topic                       string
	Broker                      broker.Broker
	ProcessedMessagesRepository repositories.ProcessedMessagesRepository
	Tracer                      trace.Tracer
	AppLogger                   *logger.AppLogger
	EmailClient                 *email.Client
	EmailRepo                   repositories.EmailNotificationTargetRepository
}

func NewKafkaListener(params NewKafkaListenerParams) *kafkaListener {
	return &kafkaListener{
		broker:                      params.Broker,
		processedMessagesRepository: params.ProcessedMessagesRepository,
		logger:                      params.AppLogger,
		topic:                       params.Topic,
		emailClient:                 params.EmailClient,
		tracer:                      params.Tracer,
		emailRepo:                   params.EmailRepo,
	}
}

func (l *kafkaListener) Start() {
	ctx := context.Background()
	l.logger.Info(ctx, "Starting Kafka listener...")

	err := l.broker.StartConsumer(l.topic, l.handleKafkaMessage)
	if err != nil {
		l.logger.Error(ctx, "Failed to start Kafka consumer", slog.Any("error", err))
	}
}

func (l *kafkaListener) handleKafkaMessage(message *broker.KafkaMessage) error {
	ctx := context.Background()

	traceparent := getHeaderValue(message.Headers, "ce-traceparent")
	otelCarrier := propagation.MapCarrier{"traceparent": traceparent}
	otel.GetTextMapPropagator().Inject(ctx, otelCarrier)

	ctx, span := l.tracer.Start(ctx, "kafkaListener.handleKafkaMessage")
	defer span.End()

	ceId := getHeaderValue(message.Headers, "ce-id")
	ceType := getHeaderValue(message.Headers, "ce-type")
	ceSource := getHeaderValue(message.Headers, "ce-source")

	span.SetAttributes(
		attribute.String("ce-id", ceId),
		attribute.String("ce-type", ceType),
		attribute.String("ce-source", ceSource),
	)

	switch ceType {
	case DonationStatusChangedEventType:
		err := l.handleDonationStatusChanged(ctx, message)
		if err != nil {
			span.SetStatus(codes.Error, "failed to handle donation status changed")
			span.RecordError(err)
			l.logger.Error(ctx, "Failed to handle DonationStatusChangedEvent", slog.Any("error", err))
		}
		return err
	case NewDonationCommittedEventType:
		err := l.handleNewDonationCommitted(ctx, message)
		if err != nil {
			span.SetStatus(codes.Error, "failed to handle new donation committed")
			span.RecordError(err)
			l.logger.Error(ctx, "Failed to handle NewDonationCommittedEvent", slog.Any("error", err))
		}
		return err
	default:
		l.logger.Warn(ctx, "Received unsupported event type", slog.String("ce-type", ceType))
		span.SetStatus(codes.Ok, "nothing to do")
		return nil
	}
}

func getHeaderValue(headers []kafka.Header, key string) string {
	for _, h := range headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (l *kafkaListener) broadcastEmailNotification(ctx context.Context, subject, body string) error {
	ctx, span := l.tracer.Start(ctx, "broadcastEmailNotification")
	defer span.End()

	targets, err := l.emailRepo.List(ctx)
	if err != nil {
		l.logger.Error(ctx, "Failed to list email notification targets", slog.Any("error", err))
		span.SetStatus(codes.Error, "Failed to list email notification targets")
		span.RecordError(err)
		return err
	}

	recipients := make([]string, 0, len(targets))
	for _, target := range targets {
		recipients = append(recipients, target.Email)
	}

	if len(recipients) == 0 {
		l.logger.Warn(ctx, "No email notification targets configured, skipping notification")
		span.SetStatus(codes.Ok, "no recipients")
		return nil
	}

	err = l.emailClient.SendEmail(ctx, recipients, subject, body)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to send email notification")
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}
