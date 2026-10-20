package homework20_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	memcacheAPI "go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/api"
	memcacheStore "go-core-4/lessons/20-NoSQL/homework-20/memcache/pkg/store"
	shortenerAPI "go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/api"
	shortenerClient "go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/cacheclient"
	shortenerStorage "go-core-4/lessons/20-NoSQL/homework-20/shortener/pkg/storage"
)

// TestE2E_Lynks_FullFlow тестирует полный цикл работы архитектуры Lynks:
// 1. Запуск микросервиса кэширования memcache.
// 2. Запуск микросервиса shortener с подключением к memcache.
// 3. Создание короткой ссылки через POST /.
// 4. Первый запрос (Cache Miss) -> обращение к СУБД -> асинхронный прогрев кэша -> редирект (HTTP 302).
// 5. Проверка наличия пары в memcache через GET /api/v1/cache/{short}.
// 6. Второй запрос (Cache Hit) -> мгновенный редирект из кэша (HTTP 302).
func TestE2E_Lynks_FullFlow(t *testing.T) {
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	// 1. Поднимаем тестовый HTTP-сервер Memcache
	cacheStore := memcacheStore.NewMemoryStore(24 * time.Hour)
	defer cacheStore.Close()
	memAPI := memcacheAPI.New(cacheStore, logger)
	memServer := httptest.NewServer(memAPI.Router())
	defer memServer.Close()

	// 2. Поднимаем сервис Shortener
	mainStorage := shortenerStorage.NewMemoryStorage()
	defer mainStorage.Close()
	cacheHTTPClient := shortenerClient.NewHTTPClient(memServer.URL, 2*time.Second)
	shortAPI := shortenerAPI.New(mainStorage, cacheHTTPClient, "https://lynks.org", logger)
	shortServer := httptest.NewServer(shortAPI.Router())
	defer shortServer.Close()

	// 3. Создание ссылки: POST /
	targetURL := "https://github.com/thinknetica/go_course_4"
	createBody, _ := json.Marshal(map[string]string{"destination": targetURL})
	resp, err := http.Post(shortServer.URL+"/", "application/json", bytes.NewReader(createBody))
	if err != nil {
		t.Fatalf("failed to POST link: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	var created shortenerAPI.ShortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.HasPrefix(created.ShortURL, "https://lynks.org/") || created.Destination != targetURL {
		t.Fatalf("unexpected created response: %+v", created)
	}

	shortID := strings.TrimPrefix(created.ShortURL, "https://lynks.org/")

	// 4. Первый переход: GET /{short} (Cache Miss -> Storage -> Cache Set -> Redirect)
	noRedirectClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	getResp, err := noRedirectClient.Get(shortServer.URL + "/" + shortID)
	if err != nil {
		t.Fatalf("failed to GET link: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", getResp.StatusCode)
	}
	if getResp.Header.Get("Location") != targetURL {
		t.Errorf("expected Location %s, got %s", targetURL, getResp.Header.Get("Location"))
	}

	// 5. Ждем прогрев кэша в memcache
	var cachedVal string
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		val, err := cacheStore.Get(context.Background(), shortID)
		if err == nil && val == targetURL {
			cachedVal = val
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if cachedVal != targetURL {
		t.Fatalf("cache was not populated in memcache service")
	}

	// 6. Проверяем API самого memcache: GET /api/v1/cache/{short}
	cacheResp, err := http.Get(memServer.URL + "/api/v1/cache/" + shortID)
	if err != nil {
		t.Fatalf("failed to call memcache API: %v", err)
	}
	defer cacheResp.Body.Close()

	if cacheResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from memcache API, got %d", cacheResp.StatusCode)
	}

	var pair memcacheAPI.CachePair
	_ = json.NewDecoder(cacheResp.Body).Decode(&pair)
	if pair.Destination != targetURL || pair.ShortURL != shortID {
		t.Errorf("unexpected cache pair from memcache API: %+v", pair)
	}

	// 7. Второй переход: GET /{short} (Cache Hit -> прямой быстрый редирект из кэша)
	getResp2, err := noRedirectClient.Get(shortServer.URL + "/" + shortID)
	if err != nil {
		t.Fatalf("failed on second GET: %v", err)
	}
	defer getResp2.Body.Close()

	if getResp2.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found on Cache Hit, got %d", getResp2.StatusCode)
	}
	if getResp2.Header.Get("Location") != targetURL {
		t.Errorf("expected Location %s, got %s", targetURL, getResp2.Header.Get("Location"))
	}

	// 8. Проверка Prometheus метрик на обоих сервисах
	mResp1, err := http.Get(shortServer.URL + "/metrics")
	if err != nil {
		t.Fatalf("failed to get shortener metrics: %v", err)
	}
	defer mResp1.Body.Close()
	mBody1, _ := io.ReadAll(mResp1.Body)
	if !strings.Contains(string(mBody1), "cache_hits_total") || !strings.Contains(string(mBody1), "cache_misses_total") {
		t.Errorf("missing cache metrics in shortener")
	}

	mResp2, err := http.Get(memServer.URL + "/metrics")
	if err != nil {
		t.Fatalf("failed to get memcache metrics: %v", err)
	}
	defer mResp2.Body.Close()
	mBody2, _ := io.ReadAll(mResp2.Body)
	if !strings.Contains(string(mBody2), "http_requests_total") {
		t.Errorf("missing http_requests_total in memcache")
	}
}
