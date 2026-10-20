// Package config реализует управление конфигурацией микросервиса (слайды 12–13 лекции №18).
// Поддерживает считывание настроек из переменных окружения с разумными значениями по умолчанию.
package config

import (
	"os"
)

// Config содержит параметры работы микросервиса.
type Config struct {
	Port     string // Порт HTTP-сервера (например, ":8080")
	BaseURL  string // Базовый URL для формирования коротких ссылок
	LogLevel string // Уровень логирования
}

// Load считывает конфигурацию из окружения.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost" + port
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return Config{
		Port:     port,
		BaseURL:  baseURL,
		LogLevel: logLevel,
	}
}
