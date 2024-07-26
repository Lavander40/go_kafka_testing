package config

import (
	"flag"
	"fmt"
	"os"
)


// Config type stores all configurational information for other modules 
type Config struct {
	LogLevel       string
	ApiPort        string
	PostgreConnect string
	KafkaConnect   string
}

// MustLoad returns formed configuration for the application or stops completely it with fatal error 
func MustLoad() *Config {
	var c Config
	var pUser, pPassword, pHost, kHost string

	// reading provided flags
	flag.StringVar(&c.LogLevel, "log", "debug", "set level for logger")
	flag.StringVar(&pUser, "user", "", "user for postgre")
	flag.StringVar(&pPassword, "pass", "", "password for postgre")
	flag.StringVar(&pHost, "phost", "", "host for postgre")
	flag.StringVar(&kHost, "khost", "", "host for kafka")
	flag.Parse()

	// if flags are empty then read from env variables
	if c.LogLevel == "" {
		c.LogLevel = os.Getenv("LOG_LEVEL")
	}
	if pUser == "" {
		pUser = os.Getenv("POSTGRES_USER")
	}
	if pPassword == "" {
		pPassword = os.Getenv("POSTGRES_PASSWORD")
	}
	if pHost == "" {
		pHost = os.Getenv("POSTGRES_HOST")
	}
	if kHost == "" {
		kHost = os.Getenv("KAFKA_HOST")
	}

	// error on required parameters
	if pUser == "" || pPassword == "" {
		panic("No info for postgre database connect provided (-user and -pass flags or env variables are required)")
	}

	c.ApiPort = "8090"
	c.PostgreConnect = fmt.Sprintf("postgres://%s:%s@%s:5432/messages?sslmode=disable", pUser, pPassword, pHost)
	c.KafkaConnect = fmt.Sprintf("%s:9092", kHost)

	return &c
}
