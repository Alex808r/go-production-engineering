package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	// ErrNotFound возвращается, если ключ отсутствует в кэше или истек его TTL.
	ErrNotFound = errors.New("ключ не найден в кэше")
	// ErrEmptyKey возвращается при попытке обращения с пустым ключом.
	ErrEmptyKey = errors.New("ключ не может быть пустым")
)

// DefaultTTL составляет 24 часа в соответствии с требованиями итогового проекта Lynks.
const DefaultTTL = 24 * time.Hour

// Client представляет клиент кэширования на базе Redis.
type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

// New инициализирует новый клиент Redis с заданным временем жизни записей (TTL).
func New(addr, password string, db int, ttl time.Duration) *Client {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Client{
		rdb: rdb,
		ttl: ttl,
	}
}

// Set сохраняет пару ключ-значение в Redis с временем жизни TTL (24 часа).
func (c *Client) Set(ctx context.Context, key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}
	return c.rdb.Set(ctx, key, value, c.ttl).Err()
}

// Get возвращает значение по ключу. Если ключ отсутствует (redis.Nil), возвращает ErrNotFound.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// Delete удаляет ключ из кэша.
func (c *Client) Delete(ctx context.Context, key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	return c.rdb.Del(ctx, key).Err()
}

// Ping проверяет доступность сервера Redis.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close корректно закрывает соединения пула Redis.
func (c *Client) Close() error {
	return c.rdb.Close()
}
