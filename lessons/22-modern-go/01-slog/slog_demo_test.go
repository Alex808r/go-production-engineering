package slogdemo

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestLogUserLogin_Masking(t *testing.T) {
	u := User{
		ID:       42,
		Email:    "alice@example.com",
		Password: "super-secret-password-123",
	}

	output := LogToBuffer(u, "192.168.1.1")

	// 1. Убеждаемся, что исходный пароль не попал в лог
	if strings.Contains(output, "super-secret-password-123") {
		t.Fatalf("критическая ошибка безопасности: секретный пароль обнаружен в логах: %s", output)
	}

	// 2. Убеждаемся, что значение замаскировано через LogValuer
	if !strings.Contains(output, "******") {
		t.Fatalf("ожидалось маскированное значение '******', лог: %s", output)
	}

	// 3. Проверяем валидность JSON
	var parsed map[string]any
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("невалидный JSON в логе: %v", err)
	}

	if parsed["msg"] != "user authenticated" {
		t.Errorf("msg mismatch: %v", parsed["msg"])
	}
	if parsed["remote_ip"] != "192.168.1.1" {
		t.Errorf("remote_ip mismatch: %v", parsed["remote_ip"])
	}
}

func BenchmarkSlog_JSONHandler(b *testing.B) {
	logger := NewJSONLogger(io.Discard, slog.LevelInfo)
	ctx := context.Background()
	u := User{ID: 10, Email: "bench@test.com", Password: "secret"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "bench event", slog.Any("user", u))
	}
}

func BenchmarkSlog_TextHandler(b *testing.B) {
	logger := NewTextLogger(io.Discard, slog.LevelInfo)
	ctx := context.Background()
	u := User{ID: 10, Email: "bench@test.com", Password: "secret"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "bench event", slog.Any("user", u))
	}
}
