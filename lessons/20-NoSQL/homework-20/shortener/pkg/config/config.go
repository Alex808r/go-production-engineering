// Package config загружает параметры конфигурации микросервиса shortener.
package config

import (
	"os"
	"strings"
)

// Config содержит конфигурационные переменные для shortener.
type Config struct {
	Port        string
	BaseURL     string
	PgDSN       string
	MemcacheURL string
}

// Load считывает переменные окружения.
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
		baseURL = "https://lynks.org"
	}

	pgDSN := os.Getenv("PG_DSN")
	if pgDSN == "" {
		pgDSN = "postgres://postgres:postgres@localhost:5432/lynks?sslmode=disable"
	}

	memcacheURL := os.Getenv("MEMCACHE_URL")
	if memcacheURL == "" {
		memcacheURL = "http://localhost:8082"
	}

	return Config{
		Port:        port,
		BaseURL:     baseURL,
		PgDSN:       pgDSN,
		MemcacheURL: memcacheURL,
	}
}
