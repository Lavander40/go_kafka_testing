package main

import (
	"context"
	"go_kafka_testing/internal/config"
	"go_kafka_testing/internal/kafka"
	"go_kafka_testing/internal/server"
	"go_kafka_testing/internal/storage"
	"log/slog"
	"os"
)

const (
	envLocal = "debug"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	log := setupLogger(cfg.LogLevel)
	if log == nil {
		panic("non existing log level was given. Set flag to debug|dev|prod")
	}
	log.Debug("logger init done")
	log.Debug("config state: ", slog.Any("config", cfg))

	s := storage.New(ctx, cfg.PostgreConnect)
	if err := s.Connect(); err != nil {
		log.Error("Failed to connect to storage: ", slog.Any("error", err))
		panic(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Error("Error closing storage: ", slog.Any("error", err))
			panic(err)
		}
	}()
	log.Debug("storage init done")

	var k kafka.KafkaWriter = kafka.New(ctx, cfg.KafkaConnect, "messages")
	if err := k.InitWriter(); err != nil {
		log.Error("Failed to connect to kafka: ", slog.Any("error", err))
		panic(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Error("Error closing kafka: ", slog.Any("error", err))
			panic(err)
		}
	}()
	log.Debug("kafka init done")
	
	server := server.New(log, cfg.ApiPort, s, k)
	if err := server.Start(); err != nil {
		log.Error("error during app launch: ", slog.Any("error", err))
		panic(err)
	}
	log.Info("Starting api server")
}

func setupLogger(level string) (log *slog.Logger) {
	switch level {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return
}
