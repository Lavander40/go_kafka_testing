package kafka

import (
	"context"
	"go_kafka_testing/internal/domain"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter interface {
	InitWriter() error
	SendMessage(*domain.Message) error
}

type KafkaReader interface {
	InitReader() error
	ReadMessage() (*domain.Message, error)
}

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

func (k *Kafka) InitWriter() error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(k.kafkaConnect),
		Topic:    k.topic,
		Balancer: &kafka.LeastBytes{},
	}

	k.writer = writer
	return nil
}

func (k *Kafka) InitReader() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.kafkaConnect},
		Topic:   k.topic,
		GroupID: k.topic + "-processor",
	})

	k.reader = reader
	return nil
}

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
