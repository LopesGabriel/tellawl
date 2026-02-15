package broker

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"go.opentelemetry.io/otel/log/noop"
)

func TestKafkaBroker(t *testing.T) {
	appLogger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
		LoggerProvider: noop.NewLoggerProvider(),
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Shutdown(t.Context())

	t.Run("Start new Producer", func(t *testing.T) {
		kafkaBroker, err := NewKafkaBroker(NewKafkaBrokerArgs{
			BootstrapServers: []string{"localhost:29092"},
			Topic:            "example_topic",
			Service:          "test-runner",
			Logger:           appLogger,
		})
		if err != nil {
			t.Errorf("failed to start broker: %v", err)
			return
		}
		defer kafkaBroker.Close()

		payload, err := json.Marshal(map[string]any{
			"firstName": "Gabriel",
			"lastName":  "Lopes",
			"email":     "lopesgabriel0199@gmail.com",
			"id":        "1578953",
		})
		if err != nil {
			t.Fatalf("failed to marshal payload")
			return
		}

		err = kafkaBroker.Produce(t.Context(), &Message{
			EventType: "dev.lopesgabriel.member-service.member.created",
			Key:       "example-key",
			Payload:   payload,
		})
		if err != nil {
			t.Error(err)
		}

		wait := make(chan bool)

		kafkaBroker.StartConsumer("example_topic", func(msg *kafka.Message) error {
			appLogger.Info(
				t.Context(),
				"New message received by callback function",
				slog.Any("Value", msg.Value),
				slog.Any("Key", msg.Key),
				slog.Any("Headers", msg.Headers),
			)

			wait <- true
			return nil
		})

		<-wait
	})
}
