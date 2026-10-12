package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore реализует Store с использованием СУБД Redis.
type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisStore создает новое подключение к Redis.
func NewRedisStore(addr, password string, db int, ttl time.Duration) (*RedisStore, error) {
	if ttl <= 0 || ttl > 24*time.Hour {
		ttl = 24 * time.Hour // По ТЗ: не более 24 часов
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis (%s): %w", addr, err)
	}

	return &RedisStore{
		client: client,
		ttl:    ttl,
	}, nil
}

// Set сохраняет пару (shortURL -> destination) в Redis с временем жизни TTL (24 часа).
func (r *RedisStore) Set(ctx context.Context, shortURL, destination string) error {
	key := strings.TrimSpace(shortURL)
	if key == "" {
		return ErrEmptyKey
	}

	if err := r.client.Set(ctx, "lynks:"+key, destination, r.ttl).Err(); err != nil {
		return fmt.Errorf("ошибка записи в Redis: %w", err)
	}
	return nil
}

// Get извлекает destination по ключу короткой ссылки.
func (r *RedisStore) Get(ctx context.Context, shortURL string) (string, error) {
	key := strings.TrimSpace(shortURL)
	if key == "" {
		return "", ErrEmptyKey
	}

	val, err := r.client.Get(ctx, "lynks:"+key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("ошибка чтения из Redis: %w", err)
	}

	return val, nil
}

// Delete удаляет ключ из Redis.
func (r *RedisStore) Delete(ctx context.Context, shortURL string) error {
	key := strings.TrimSpace(shortURL)
	if key == "" {
		return ErrEmptyKey
	}
	return r.client.Del(ctx, "lynks:"+key).Err()
}

// Close закрывает пул соединений с Redis.
func (r *RedisStore) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
