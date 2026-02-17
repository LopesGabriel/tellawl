package models

import (
	"time"

	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/events"
)

type ProcessedMessage struct {
	MessageID   string
	Recepient   string
	Subject     string
	Sender      string
	Body        string
	ProcessedAt time.Time
	events      []events.DomainEvent
}

func NewProcessedMessage(messageID, recepient, subject, sender, body string) *ProcessedMessage {
	event := events.EmailReceivedEvent{
		MessageId: messageID,
		Recipient: recepient,
		Sender:    sender,
		Subject:   subject,
		Body:      body,
		Timestamp: time.Now(),
	}

	return &ProcessedMessage{
		MessageID:   messageID,
		Recepient:   recepient,
		Subject:     subject,
		Sender:      sender,
		Body:        body,
		ProcessedAt: time.Now(),
		events:      []events.DomainEvent{event},
	}
}

func (m *ProcessedMessage) GetEvents() []events.DomainEvent {
	return m.events
}

func (m *ProcessedMessage) ClearEvents() {
	m.events = []events.DomainEvent{}
}
