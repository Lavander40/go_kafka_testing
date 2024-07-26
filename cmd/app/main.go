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
	// loading configuration
	cfg := config.MustLoad()
	// creating context
	ctx := context.Background()
	// setuping logger to work with configured level
	log := setupLogger(cfg.LogLevel)
	if log == nil {
		panic("non existing log level was given. Set flag to debug|dev|prod")
	}
	log.Debug("logger init done")
	log.Debug("config state: ", slog.Any("config", cfg))

	// initialization of storage and starting up connection
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

	// passing kafka configuration and setuping writer
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

	// creating and initializating api server
	server := server.New(log, cfg.ApiPort, s, k)
	if err := server.Start(); err != nil {
		log.Error("error during app launch: ", slog.Any("error", err))
		panic(err)
	}
	log.Info("Starting api server")
}

// setupLogger configures logger to work on set level.
// Logger handles work with debug/dev/prod levels.
func setupLogger(level string) (log *slog.Logger) {
	switch level {
	case envLocal:
		// Debug - uses text format with messages from DEBUG and higher levels;
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		// Dev - uses JSON output with messages from DEBUG and higher;
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		// Prod - uses JSON output but with messages from INFO level and higher.
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return
}
