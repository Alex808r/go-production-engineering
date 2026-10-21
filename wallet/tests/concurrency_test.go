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

func getTestConfig() *config.Config {
	cfg, err := config.Load("../config.env")
	if err != nil {
		cfg = &config.Config{}
	}

	// Use local test database if available
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

	// Create wallet
	_, err := repo.Create(ctx, walletID, initialBalance)
	if err != nil {
		t.Fatalf("failed to create initial wallet: %v", err)
	}

	const totalRequests = 1000
	const depositCount = 500
	const withdrawCount = 500
	depositAmount := decimal.NewFromInt(10)
	withdrawAmount := decimal.NewFromInt(5)

	var success200Count int64
	var error500Count int64
	var otherErrorCount int64

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 200,
			MaxConnsPerHost:     200,
		},
	}

	var wg sync.WaitGroup
	startBarrier := make(chan struct{})

	// Launch 500 deposits and 500 withdrawals simultaneously
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
				"valletId":      walletID.String(),
				"operationType": opType,
				"amount":        amt,
			}
			bodyBytes, _ := json.Marshal(payload)

			<-startBarrier // synchronize start of all 1000 requests

			resp, err := client.Post(ts.URL+"/api/v1/wallet", "application/json", bytes.NewReader(bodyBytes))
			if err != nil {
				atomic.AddInt64(&error500Count, 1)
				t.Logf("client error: %v", err)
				return
			}
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, resp.Body)

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

	// Trigger all 1000 requests at the exact same moment
	start := time.Now()
	close(startBarrier)
	wg.Wait()
	duration := time.Since(start)

	t.Logf("1000 requests completed in %v (throughput: %.2f req/sec)",
		duration, float64(totalRequests)/duration.Seconds())
	t.Logf("Results: 200 OK = %d, 50x Errors = %d, Other Errors = %d",
		success200Count, error500Count, otherErrorCount)

	// Requirement: "Ни один запрос не должен быть не обработан (50Х error)"
	if error500Count != 0 {
		t.Fatalf("CRITICAL: Received %d 500 Internal Server Errors under load!", error500Count)
	}

	if success200Count != totalRequests {
		t.Fatalf("expected all %d requests to succeed with 200 OK, got: %d", totalRequests, success200Count)
	}

	// Verify balance integrity:
	// initial (100000) + 500 * 10 - 500 * 5 = 100000 + 5000 - 2500 = 102500
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

func TestConcurrency_WithdrawExceedingBalance(t *testing.T) {
	ts, repo, _ := setupTestApp(t)
	if ts == nil {
		return
	}
	defer ts.Close()
	defer repo.Close()

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := decimal.NewFromInt(50) // only 50

	_, err := repo.Create(ctx, walletID, initialBalance)
	if err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}

	// Try 10 concurrent withdrawals of 10 each
	// Exactly 5 must succeed (5 * 10 = 50), and 5 must fail with 400 Bad Request
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

	// Final balance must be 0
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
