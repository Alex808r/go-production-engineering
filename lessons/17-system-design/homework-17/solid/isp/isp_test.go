package isp_test

import (
	"context"
	"testing"

	"go-core-4/lessons/17-system-design/homework-17/solid/isp"
)

// MockReader для тестирования ReadOnlyService реализует только 1 метод, а не весь CRUD!
type MockReader struct {
	item isp.Item
	err  error
}

func (m *MockReader) Get(_ context.Context, _ int64) (isp.Item, error) {
	return m.item, m.err
}

func TestReadOnlyService_ISP(t *testing.T) {
	ctx := context.Background()

	// 1. Использование минимального мока — подтверждение удобства ISP
	mock := &MockReader{item: isp.Item{ID: 42, Name: "Clean Architecture"}}
	svc := isp.NewReadOnlyService(mock)

	name, err := svc.GetItemName(ctx, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "Clean Architecture" {
		t.Errorf("expected 'Clean Architecture', got %q", name)
	}

	// 2. Использование реального хранилища через узкую проекцию Reader
	realStorage := isp.NewMemStorage()
	_ = realStorage.Save(ctx, isp.Item{ID: 100, Name: "Go Programming"})

	svcReal := isp.NewReadOnlyService(realStorage)
	realName, err := svcReal.GetItemName(ctx, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if realName != "Go Programming" {
		t.Errorf("expected 'Go Programming', got %q", realName)
	}
}
