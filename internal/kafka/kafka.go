package kafka

import (
	"context"
	"go_kafka_testing/internal/domain"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter interface contains functions to connect to kafka for sending messages to the topic
type KafkaWriter interface {
	InitWriter() error
	SendMessage(*domain.Message) error
}

// KafkaReader interface contains functions to connect to kafka for receiving messages from the topic
type KafkaReader interface {
	InitReader() error
	ReadMessage() (*domain.Message, error)
}

// Kafka type stores information about kafka broker and opened connections for writing and reading messages from it.
type Kafka struct {
	ctx          context.Context
	writer       *kafka.Writer
	reader       *kafka.Reader
	kafkaConnect string
	topic        string
}

func New(ctx context.Context, kafkaConnect string, topic string) *Kafka {
	return &Kafka{
		ctx:          ctx,
		kafkaConnect: kafkaConnect,
		topic:        topic,
	}
}

// InitWriter creates connect for sending messages to kafka
func (k *Kafka) InitWriter() error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(k.kafkaConnect),
		Topic:    k.topic,
		Balancer: &kafka.LeastBytes{},
	}

	k.writer = writer
	return nil
}

// InitReader creates connect to read messages from kafka
func (k *Kafka) InitReader() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.kafkaConnect},
		Topic:   k.topic,
		GroupID: k.topic + "-processor",
	})

	k.reader = reader
	return nil
}

// Close fucn closes all existing connections
func (k *Kafka) Close() error {
	if k.writer != nil {
		if err := k.writer.Close(); err != nil {
			return err
		}
	}
	if k.reader != nil {
		if err := k.reader.Close(); err != nil {
			return err
		}
	}
	return nil
}
