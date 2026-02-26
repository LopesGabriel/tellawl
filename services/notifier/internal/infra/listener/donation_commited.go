package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/lopesgabriel/tellawl/packages/broker"
	"go.opentelemetry.io/otel/codes"
)

const NewDonationCommittedEventType = "dev.lopesgabriel.casanova.donation.committed"

// NewDonationCommittedEvent is emitted when a user commits to donating a desired item.
type NewDonationCommittedEvent struct {
	DonationID      string    `json:"donation_id"`
	DesiredItemID   string    `json:"desired_item_id"`
	DesiredItemName string    `json:"desired_item_name"`
	DonorID         string    `json:"donor_id"`
	DonorName       string    `json:"donor_name"`
	Amount          float64   `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
}

func (e NewDonationCommittedEvent) AggregateID() string {
	return e.DonationID
}

func (e NewDonationCommittedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

func (l *kafkaListener) handleNewDonationCommitted(ctx context.Context, message *broker.KafkaMessage) error {
	ctx, span := l.tracer.Start(ctx, "handleNewDonationCommitted")
	defer span.End()

	var event NewDonationCommittedEvent
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		l.logger.Error(ctx, "Failed to unmarshal NewDonationCommittedEvent", slog.Any("error", err))
		span.SetStatus(codes.Error, "Failed to unmarshal event")
		span.RecordError(err)
		return err
	}

	subject := fmt.Sprintf("Nova doação: %s comprometeu-se com '%s'", event.DonorName, event.DesiredItemName)
	body := fmt.Sprintf("O doador '%s' comprometeu-se a doar R$ %.2f para o item '%s'.", event.DonorName, event.Amount, event.DesiredItemName)
	err = l.broadcastEmailNotification(ctx, subject, body)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to broadcast email notification")
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}
