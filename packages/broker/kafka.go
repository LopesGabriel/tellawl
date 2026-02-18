package broker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"go.opentelemetry.io/otel/trace"
)

type kafkaBroker struct {
	closeChan       chan struct{}
	doneChan        chan struct{}
	consumerStarted bool
	consumer        *kafka.Consumer
	logger          *logger.AppLogger
	topic           string
	producer        *kafka.Producer
	serviceName     string
}

// NewKafkaBrokerArgs holds the arguments for creating a new Kafka broker.
type NewKafkaBrokerArgs struct {
	BootstrapServers []string
	Service          string
	Topic            string
	Logger           *logger.AppLogger
}

// NewKafkaBroker creates a new Kafka broker with a producer and consumer.
func NewKafkaBroker(args NewKafkaBrokerArgs) (Broker, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(args.BootstrapServers, ","),
		"allow.auto.create.topics": "true",
	})
	if err != nil {
		return nil, err
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(args.BootstrapServers, ","),
		"group.id":                 fmt.Sprintf("%s-group", args.Service),
		"auto.offset.reset":        "earliest",
		"allow.auto.create.topics": "true",
		"session.timeout.ms":       6000,
		"max.poll.interval.ms":     300000, // 5 minutos
	})
	if err != nil {
		return nil, err
	}

	return &kafkaBroker{
		topic:       args.Topic,
		serviceName: args.Service,
		consumer:    c,
		producer:    p,
		closeChan:   make(chan struct{}),
		doneChan:    make(chan struct{}),
		logger:      args.Logger,
	}, nil
}

func (k *kafkaBroker) Close() error {
	k.producer.Close()
	if k.consumerStarted {
		close(k.closeChan)
		<-k.doneChan
	}
	return k.consumer.Close()
}

func (k *kafkaBroker) Produce(ctx context.Context, message *Message) error {
	msg := &kafka.Message{
		Value: message.Payload,
		Key:   []byte(message.Key),
		Headers: []kafka.Header{
			{Key: "ce-id", Value: []byte(message.EventId)},
			{Key: "ce-type", Value: []byte(message.EventType)},
			{Key: "ce-source", Value: []byte(fmt.Sprintf("dev.lopesgabriel.tellawl.%s", k.serviceName))},
			{Key: "ce-traceparent", Value: []byte(buildTraceparent(ctx))},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: -1,
		},
	}

	wait := make(chan kafka.Event)
	err := k.producer.Produce(msg, wait)
	if err != nil {
		return err
	}

	<-wait
	close(wait)
	return nil
}

func (k *kafkaBroker) StartConsumer(topic string, callback CallbackFunction) error {
	err := k.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return err
	}

	k.consumerStarted = true

	go func() {
		for {
			select {
			case <-k.closeChan:
				k.logger.Debug(context.TODO(), "Closing Consumer")
				k.logger.Debug(context.TODO(), "consumer stopped")
				close(k.doneChan)
				return
			default:
				msg, err := k.consumer.ReadMessage(time.Second)
				if err != nil {
					if !err.(kafka.Error).IsTimeout() {
						k.logger.Error(context.TODO(), "consumer error", slog.Any("error", err))
					}

					k.logger.Debug(context.TODO(), "No messages for the previous second")
					continue
				}

				err = callback(msg)
				if err == nil {
					k.logger.Debug(context.TODO(), "committing message", slog.Any("Key", msg.Key))
					_, err := k.consumer.CommitMessage(msg)
					if err != nil {
						k.logger.Error(context.TODO(), "failed to commit message", slog.Any("Key", msg.Key))
					}
				}
			}
		}
	}()

	return nil
}

func buildTraceparent(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	if !sc.IsValid() {
		return ""
	}

	traceID := sc.TraceID().String()
	spanID := sc.SpanID().String()

	flags := "00"
	if sc.IsSampled() {
		flags = "01"
	}

	return "00-" + traceID + "-" + spanID + "-" + flags
}
