package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_FromEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, "config.env")

	content := `
# Test Comment
APP_PORT=9090
APP_ENV=test
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=test_user
DB_PASSWORD=secret
DB_NAME=test_db
DB_SSLMODE=require
DB_MAX_CONNS=75
DB_MIN_CONNS=15
DB_MAX_IDLE_TIME=2m
DB_MAX_CONN_LIFETIME=30m
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=20s
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test env: %v", err)
	}

	cfg, err := Load(envPath)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.AppPort != "9090" {
		t.Errorf("expected AppPort 9090, got %s", cfg.AppPort)
	}
	if cfg.DBHost != "127.0.0.1" {
		t.Errorf("expected DBHost 127.0.0.1, got %s", cfg.DBHost)
	}
	if cfg.DBMaxConns != 75 {
		t.Errorf("expected DBMaxConns 75, got %d", cfg.DBMaxConns)
	}
	if cfg.DBMaxIdleTime != 2*time.Minute {
		t.Errorf("expected DBMaxIdleTime 2m, got %v", cfg.DBMaxIdleTime)
	}
	if cfg.DSN() != "postgres://test_user:secret@127.0.0.1:5432/test_db?sslmode=require" {
		t.Errorf("unexpected DSN: %s", cfg.DSN())
	}
}
