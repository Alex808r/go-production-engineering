// Package cacheclient реализует HTTP-клиент к микросервису кэширования memcache.
package cacheclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	// ErrCacheMiss возвращается при отсутствии ключа в кэше.
	ErrCacheMiss = errors.New("ключ отсутствует в кэше")
)

// Client интерфейс взаимодействия с сервисом кэша.
type Client interface {
	Get(ctx context.Context, shortURL string) (string, error)
	Set(ctx context.Context, shortURL, destination string) error
}

// HTTPClient реализует Client через HTTP-запросы к микросервису memcache.
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPClient создает новый HTTP-клиент кэша.
func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	baseURL = strings.TrimRight(baseURL, "/")
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type cachePayload struct {
	ShortURL    string `json:"shortUrl"`
	Destination string `json:"destination"`
}

// Get запрашивает у сервиса memcache оригинальный URL по короткой ссылке.
func (c *HTTPClient) Get(ctx context.Context, shortURL string) (string, error) {
	if c == nil || c.baseURL == "" {
		return "", ErrCacheMiss
	}

	url := fmt.Sprintf("%s/api/v1/cache/%s", c.baseURL, shortURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка сетевого запроса к кэшу: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrCacheMiss
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("кэш вернул статус %d", resp.StatusCode)
	}

	var payload cachePayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("ошибка парсинга ответа кэша: %w", err)
	}

	return payload.Destination, nil
}

// Set отправляет пару (shortURL -> destination) в сервис memcache.
func (c *HTTPClient) Set(ctx context.Context, shortURL, destination string) error {
	if c == nil || c.baseURL == "" {
		return nil
	}

	bodyBytes, err := json.Marshal(cachePayload{
		ShortURL:    shortURL,
		Destination: destination,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/cache", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки в кэш: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("кэш вернул ошибку: статус %d", resp.StatusCode)
	}

	return nil
}
