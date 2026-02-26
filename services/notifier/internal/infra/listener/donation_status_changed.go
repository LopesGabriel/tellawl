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

const DonationStatusChangedEventType = "dev.lopesgabriel.casanova.donation.status_changed"

// DonationStatus represents the lifecycle state of a donation.
type DonationStatus string

const (
	DonationStatusPending  DonationStatus = "pending"
	DonationStatusPaid     DonationStatus = "paid"
	DonationStatusCanceled DonationStatus = "canceled"
)

// DonationStatusChangedEvent is emitted when a donation's status transitions.
type DonationStatusChangedEvent struct {
	DonationID      string         `json:"donation_id"`
	DesiredItemID   string         `json:"desired_item_id"`
	DesiredItemName string         `json:"desired_item_name"`
	DonorName       string         `json:"donor_name"`
	OldStatus       DonationStatus `json:"old_status"`
	NewStatus       DonationStatus `json:"new_status"`
	Timestamp       time.Time      `json:"timestamp"`
}

func (e DonationStatusChangedEvent) AggregateID() string {
	return e.DonationID
}

func (e DonationStatusChangedEvent) OccurredAt() time.Time {
	return e.Timestamp
}

func (l *kafkaListener) handleDonationStatusChanged(ctx context.Context, message *broker.KafkaMessage) error {
	ctx, span := l.tracer.Start(ctx, "handleDonationStatusChanged")
	defer span.End()

	var event DonationStatusChangedEvent
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		l.logger.Error(ctx, "Failed to unmarshal DonationStatusChangedEvent", slog.Any("error", err))
		span.SetStatus(codes.Error, "Failed to unmarshal event")
		span.RecordError(err)
		return err
	}

	statusAntigo := "desconhecido"
	switch event.OldStatus {
	case DonationStatusPaid:
		statusAntigo = "pago"
	case DonationStatusPending:
		statusAntigo = "pendente"
	case DonationStatusCanceled:
		statusAntigo = "cancelado"
	}

	statusNovo := "desconhecido"
	switch event.NewStatus {
	case DonationStatusPaid:
		statusNovo = "pago"
	case DonationStatusPending:
		statusNovo = "pendente"
	case DonationStatusCanceled:
		statusNovo = "cancelado"
	}

	subject := fmt.Sprintf("Status de doação alterado: %s → %s", statusAntigo, statusNovo)
	body := fmt.Sprintf("Doação de '%s' para o item '%s' mudou de status de '%s' para '%s'.", event.DonorName, event.DesiredItemName, statusAntigo, statusNovo)
	err = l.broadcastEmailNotification(ctx, subject, body)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to broadcast email notification")
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}
