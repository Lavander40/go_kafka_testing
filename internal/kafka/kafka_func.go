package kafka

import (
	"encoding/json"
	"fmt"
	"go_kafka_testing/internal/domain"

	"github.com/segmentio/kafka-go"
)

// SendMessage sends formatted messages of type domain.Messae to kafka topic
func (k *Kafka) SendMessage(msg *domain.Message) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return k.writer.WriteMessages(k.ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", msg.ID)),
			Value: messageBytes,
		},
	)
}

// ReadMessage resives oldest unread message from topic
func (k *Kafka) ReadMessage() (*domain.Message, error) {
	m, err := k.reader.ReadMessage(k.ctx)
	if err != nil {
		return nil, err
	}

	msg := &domain.Message{}
	err = json.Unmarshal(m.Value, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}