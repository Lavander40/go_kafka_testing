package main

import (
	"context"
	"database/sql"
	"go_kafka_testing/internal/config"
	"go_kafka_testing/internal/domain"
	"go_kafka_testing/internal/kafka"
	"go_kafka_testing/internal/storage"
	"log"
	"time"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	s := storage.New(ctx, cfg.PostgreConnect)
	if err := s.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Fatalf("Error closing storage: %v", err)
		}
	}()

	var k kafka.KafkaReader = kafka.New(ctx, cfg.KafkaConnect, "messages")
	if err := k.InitReader(); err != nil {
		log.Fatalf("Failed to connect to kafka: %s", err)
		panic(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Fatalf("Error closing kafka: %v", err)
			panic(err)
		}
	}()
	
	messageChan := make(chan *domain.Message)
	go readMessages(k, messageChan)
	processMessages(messageChan, s)
}

func readMessages(k kafka.KafkaReader, messageChan chan<- *domain.Message) {
	for {
		msg, err := k.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		messageChan <- msg
	}
}

func processMessages(messageChan <-chan *domain.Message, s *storage.Storage) {
	for msg := range messageChan {
		log.Printf("Got message to process: %v", msg)

		time.Sleep(5 * time.Second)

		msg.Status = "processed"
		msg.ProcessedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		err := s.UpdateStatus(msg)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		log.Printf("Processed message: ID=%d, Content=%s, Status=%s, CreatedAt=%s, ProcessedAt=%v\n",
			msg.ID, msg.Content, msg.Status, msg.CreatedAt, msg.ProcessedAt)
	}
}


