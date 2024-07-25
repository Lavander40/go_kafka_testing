package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	LogLevel       string
	ApiPort        string
	PostgreConnect string
	KafkaConnect   string
}

func MustLoad() *Config {
	var c Config
	var pUser, pPassword string

	flag.StringVar(&c.LogLevel, "log", "debug", "set level for logger")
	flag.StringVar(&pUser, "user", "", "user for postgre")
	flag.StringVar(&pPassword, "pass", "", "password for postgre")
	flag.Parse()

	if c.LogLevel == "" {
		c.LogLevel = os.Getenv("LOG_LEVEL")
	}
	if pUser == "" {
		pUser = os.Getenv("POSTGRES_USER")
	}
	if pPassword == "" {
		pPassword = os.Getenv("POSTGRES_PASSWORD")
	}

	if pUser == "" || pPassword == "" {
		panic("No info for postgre database connect provided (-user and -pass flags or env variables are required)")
	}

	c.ApiPort = "8090"
	c.PostgreConnect = fmt.Sprintf("postgres://%s:%s@local:5432/messages?sslmode=disable", pUser, pPassword)
	c.KafkaConnect = "kafka:9092"

	return &c
}
