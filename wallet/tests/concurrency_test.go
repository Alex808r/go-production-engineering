package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Alex808r/wallet/config"
	httpdelivery "github.com/Alex808r/wallet/internal/delivery/http"
	"github.com/Alex808r/wallet/internal/domain"
	"github.com/Alex808r/wallet/internal/repository/postgres"
	"github.com/Alex808r/wallet/internal/service"
	internalsync "github.com/Alex808r/wallet/internal/sync"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// getTestConfig подготавливает тестовую конфигурацию для локального прогона
func getTestConfig() *config.Config {
	cfg, err := config.Load("../config.env")
	if err != nil {
		cfg = &config.Config{}
	}

	// Если запускается локально (не в CI), используем локальную тестовую базу данных
	if os.Getenv("CI") == "" && cfg.DBHost == "postgres" {
		cfg.DBHost = "127.0.0.1"
		cfg.DBName = "wallet_test"
		currentUser := os.Getenv("USER")
		if currentUser != "" {
			cfg.DBUser = currentUser
			cfg.DBPassword = ""
		}
	}
	cfg.DBMaxConns = 100
	cfg.DBMinConns = 20
	return cfg
}

// setupTestApp поднимает изолированное тестовое приложение и тестовый HTTP-сервер
func setupTestApp(t *testing.T) (*httptest.Server, *postgres.WalletRepo, service.WalletService) {
	cfg := getTestConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo, err := postgres.NewWalletRepo(ctx, cfg)
	if err != nil {
		t.Skipf("skipping integration test: cannot connect to postgres at %s: %v", cfg.DSN(), err)
		return nil, nil, nil
	}

	if err := repo.AutoMigrate(context.Background()); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	locker := internalsync.NewKeyLocker(4096)
	svc := service.NewWalletService(repo, locker)
	handler := httpdelivery.NewWalletHandler(svc)
	router := httpdelivery.NewRouter(handler)

	ts := httptest.NewServer(router)
	return ts, repo, svc
}

// TestConcurrency_1000_RPS_SingleWallet — ключевой нагрузочный стресс-тест.
// -----------------------------------------------------------------------------
// Цель теста:
// Доказать выполнение критических требований ТЗ:
// 1. Обработка 1000 одновременных запросов по ОДНОМУ кошельку.
// 2. Строго 0 ошибок 50x (Internal Server Error).
// 3. 100% сохранение финансовой консистентности (отсутствие Lost Updates).
//
// Паттерн тестирования:
// - Создаем кошелек с балансом 100 000.
// - Создаем 1000 горутин: 500 пополнений по 10 и 500 списаний по 5.
// - Используем канал-барьер `startBarrier`, чтобы выстрелить всеми 1000 запросами синхронно в одну миллисекунду.
// - Математическое ожидание: 100000 + 500*10 - 500*5 = 102500.
// -----------------------------------------------------------------------------
func TestConcurrency_1000_RPS_SingleWallet(t *testing.T) {
	ts, repo, _ := setupTestApp(t)
	if ts == nil {
		return
	}
	defer ts.Close()
	defer repo.Close()

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := decimal.NewFromInt(100000)

	// Инициализация кошелька с начальным балансом 100000
	_, err := repo.Create(ctx, walletID, initialBalance)
	if err != nil {
		t.Fatalf("failed to create initial wallet: %v", err)
	}

	const totalRequests = 1000
	const depositCount = 500
	const withdrawCount = 500
	depositAmount := decimal.NewFromInt(10)
	withdrawAmount := decimal.NewFromInt(5)

	// Атомарные счетчики статусов ответов для безопасного подсчета из 1000 горутин
	var success200Count int64
	var error500Count int64
	var otherErrorCount int64

	// Настройка HTTP-клиента с увеличенным пулом коннектов под нагрузку
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			MaxConnsPerHost:     200,
		},
	}

	var wg sync.WaitGroup
	// Канал-барьер: все горутины ждут закрытия этого канала перед отправкой HTTP-запроса
	startBarrier := make(chan struct{})

	// Запускаем 500 горутин на пополнение и 500 горутин на списание
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		isDeposit := (i < depositCount)

		go func(deposit bool) {
			defer wg.Done()

			opType := "WITHDRAW"
			amt := withdrawAmount
			if deposit {
				opType = "DEPOSIT"
				amt = depositAmount
			}

			payload := map[string]any{
				"valletId":      walletID.String(), // Проверяем работу с опечаткой из ТЗ
				"operationType": opType,
				"amount":        amt,
			}
			bodyBytes, _ := json.Marshal(payload)

			<-startBarrier // Блокировка до общего залпа

			resp, err := client.Post(ts.URL+"/api/v1/wallet", "application/json", bytes.NewReader(bodyBytes))
			if err != nil {
				atomic.AddInt64(&error500Count, 1)
				t.Logf("client error: %v", err)
				return
			}
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, resp.Body)

			// Атомарная фиксация кода ответа
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt64(&success200Count, 1)
			case http.StatusInternalServerError:
				atomic.AddInt64(&error500Count, 1)
			default:
				atomic.AddInt64(&otherErrorCount, 1)
			}
		}(isDeposit)
	}

	// ЗАЛП: закрываем канал, 1000 горутин одновременно шлют запросы
	start := time.Now()
	close(startBarrier)
	wg.Wait()
	duration := time.Since(start)

	t.Logf("1000 requests completed in %v (throughput: %.2f req/sec)",
		duration, float64(totalRequests)/duration.Seconds())
	t.Logf("Results: 200 OK = %d, 50x Errors = %d, Other Errors = %d",
		success200Count, error500Count, otherErrorCount)

	// --------------------------------------------------------------------------
	// Главная проверка ТЗ: "Ни один запрос не должен быть не обработан (50Х error)"
	// --------------------------------------------------------------------------
	if error500Count != 0 {
		t.Fatalf("CRITICAL: Received %d 500 Internal Server Errors under load!", error500Count)
	}

	if success200Count != totalRequests {
		t.Fatalf("expected all %d requests to succeed with 200 OK, got: %d", totalRequests, success200Count)
	}

	// --------------------------------------------------------------------------
	// Проверка финансовой консистентности:
	// Итоговый баланс = 100000 + (500 * 10) - (500 * 5) = 102500
	// Ни одно пополнение или списание не должно потеряться
	// --------------------------------------------------------------------------
	expectedBalance := initialBalance.Add(decimal.NewFromInt(int64(depositCount * 10))).Sub(decimal.NewFromInt(int64(withdrawCount * 5)))

	getResp, err := client.Get(ts.URL + "/api/v1/wallets/" + walletID.String())
	if err != nil {
		t.Fatalf("failed to query final balance: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from GET balance, got: %d", getResp.StatusCode)
	}

	var wResp domain.WalletResponse
	if err := json.NewDecoder(getResp.Body).Decode(&wResp); err != nil {
		t.Fatalf("failed to parse wallet response: %v", err)
	}

	if !wResp.Balance.Equal(expectedBalance) {
		t.Fatalf("CONCURRENCY DATA CORRUPTION: expected balance %v, got %v", expectedBalance, wResp.Balance)
	}

	t.Logf("SUCCESS: Final balance verified exactly: %v (zero lost updates, zero 50x errors!)", wResp.Balance)
}

// TestConcurrency_WithdrawExceedingBalance проверяет корректность списания при нехватке денег.
// При исчерпании баланса сервис обязан возвращать 400 Bad Request, а не 500 Error.
func TestConcurrency_WithdrawExceedingBalance(t *testing.T) {
	ts, repo, _ := setupTestApp(t)
	if ts == nil {
		return
	}
	defer ts.Close()
	defer repo.Close()

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := decimal.NewFromInt(50) // Всего 50 единиц

	_, err := repo.Create(ctx, walletID, initialBalance)
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}

	// Запускаем 10 одновременных списаний по 10 единиц
	// Ровно 5 должны успешно пройти (5 * 10 = 50), а остальные 5 должны отклониться с кодом 400
	const reqs = 10
	var successCount int64
	var insufficientCount int64
	var error500Count int64

	client := &http.Client{Timeout: 10 * time.Second}
	var wg sync.WaitGroup
	startBarrier := make(chan struct{})

	for i := 0; i < reqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			payload := map[string]any{
				"valletId":      walletID.String(),
				"operationType": "WITHDRAW",
				"amount":        10,
			}
			bodyBytes, _ := json.Marshal(payload)
			<-startBarrier

			resp, err := client.Post(ts.URL+"/api/v1/wallet", "application/json", bytes.NewReader(bodyBytes))
			if err != nil {
				atomic.AddInt64(&error500Count, 1)
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt64(&successCount, 1)
			case http.StatusBadRequest:
				atomic.AddInt64(&insufficientCount, 1)
			case http.StatusInternalServerError:
				atomic.AddInt64(&error500Count, 1)
			}
		}()
	}

	close(startBarrier)
	wg.Wait()

	if error500Count != 0 {
		t.Errorf("expected 0 500 errors, got: %d", error500Count)
	}
	if successCount != 5 {
		t.Errorf("expected exactly 5 successful withdrawals, got: %d", successCount)
	}
	if insufficientCount != 5 {
		t.Errorf("expected exactly 5 insufficient funds 400 errors, got: %d", insufficientCount)
	}

	// Итоговый баланс обязан быть ровно 0
	getResp, err := client.Get(fmt.Sprintf("%s/api/v1/wallets/%s", ts.URL, walletID))
	if err != nil {
		t.Fatalf("failed to get balance: %v", err)
	}
	defer getResp.Body.Close()

	var wResp domain.WalletResponse
	_ = json.NewDecoder(getResp.Body).Decode(&wResp)
	if !wResp.Balance.Equal(decimal.Zero) {
		t.Errorf("expected balance 0, got %v", wResp.Balance)
	}
}
