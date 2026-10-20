// Package config загружает конфигурацию микросервиса кэширования memcache.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config содержит настройки сервиса кэширования.
type Config struct {
	Port          string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	CacheTTL      time.Duration
}

// Load считывает переменные окружения.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8082"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}

	ttl := 24 * time.Hour
	if ttlStr := os.Getenv("CACHE_TTL"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil && parsed > 0 && parsed <= 24*time.Hour {
			ttl = parsed
		}
	}

	return Config{
		Port:          port,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
		CacheTTL:      ttl,
	}
}
