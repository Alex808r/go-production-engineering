// Package config загружает конфигурацию микросервиса аналитики.
package config

import (
	"os"
	"strings"
)

// Config содержит параметры запуска микросервиса аналитики.
type Config struct {
	Port         string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

// Load читает переменные окружения с безопасными значениями по умолчанию.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8081"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
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

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "analytics-service-group"
	}

	return Config{
		Port:         port,
		KafkaBrokers: brokers,
		KafkaTopic:   topic,
		KafkaGroupID: groupID,
	}
}
