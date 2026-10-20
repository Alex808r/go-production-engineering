package pgsql_test

import (
	"context"
	"testing"
	"time"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
	"go-core-4/lessons/16-db-apps/homework-16/pkg/db/pgsql"
)

// Проверка того, что тип *pgsql.DB статически реализует контракт db.Interface
func TestInterfaceCompliance(t *testing.T) {
	var _ db.Interface = (*pgsql.DB)(nil)
}

// Проверка корректной обработки ошибок подключения при недоступном сервере
func TestNew_ConnectionError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Заведомо недоступный адрес
	_, err := pgsql.New(ctx, "postgres://user:pass@127.0.0.1:54329/invalid_db?connect_timeout=1")
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

// Проверка пустого батча
func TestAddMovies_Empty(t *testing.T) {
	dbInstance := pgsql.NewWithPool(nil)
	err := dbInstance.AddMovies(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected nil error on empty batch, got: %v", err)
	}
}
