// Package config загружает параметры конфигурации микросервиса shortener.
package config

import (
	"os"
	"strings"
)

// Config содержит конфигурацию микросервиса URL Shortener.
type Config struct {
	Port         string
	BaseURL      string
	KafkaBrokers []string
	KafkaTopic   string
}

// Load читает переменные окружения с безопасными значениями по умолчанию.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	rawBrokers := os.Getenv("KAFKA_BROKERS")
	if rawBrokers == "" {
		rawBrokers = "localhost:29092"
	}
	brokers := strings.Split(rawBrokers, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "links-events"
	}

	return Config{
		Port:         port,
		BaseURL:      baseURL,
		KafkaBrokers: brokers,
		KafkaTopic:   topic,
	}
}
