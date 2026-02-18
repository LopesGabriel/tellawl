package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaMessage = kafka.Message
type CallbackFunction = func(*KafkaMessage) error

type Broker interface {
	Close() error
	Produce(ctx context.Context, message *Message) error
	StartConsumer(topic string, callback CallbackFunction) error
}

type Message struct {
	EventId   string
	EventType string
	Key       string
	Payload   []byte
}
