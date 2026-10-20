package traps

import (
	"context"
	"testing"
	"time"
)

func TestStructPadding(t *testing.T) {
	bad, good := GetStructSizes()
	if bad != 24 {
		t.Errorf("expected BadStruct size 24 bytes, got %d", bad)
	}
	if good != 16 {
		t.Errorf("expected GoodStruct size 16 bytes, got %d", good)
	}

	savingsPercent := float64(bad-good) / float64(bad) * 100
	t.Logf("Размер BadStruct: %d байт, GoodStruct: %d байт. Экономия: %.1f%%", bad, good, savingsPercent)
}

func TestFunctionalOptions(t *testing.T) {
	// 1. Дефолтная конфигурация
	cfg := NewServerConfig("localhost")
	if cfg.Port != 8080 || cfg.Timeout != 30*time.Second || cfg.MaxConn != 1000 {
		t.Errorf("default config mismatch: %+v", cfg)
	}

	// 2. С опциями
	cfg = NewServerConfig("0.0.0.0",
		WithPort(9090),
		WithTimeout(5*time.Second),
		WithMaxConn(500),
	)

	if cfg.Port != 9090 || cfg.Timeout != 5*time.Second || cfg.MaxConn != 500 {
		t.Errorf("custom config mismatch: %+v", cfg)
	}
}

func TestSafeSubslice(t *testing.T) {
	huge := make([]int, 1000)
	for i := range huge {
		huge[i] = i
	}

	sub := SafeSubslice(huge)
	if len(sub) != 2 || cap(sub) != 2 {
		t.Errorf("expected len 2 and cap 2 for isolated slice, got len=%d, cap=%d", len(sub), cap(sub))
	}
}

func TestProcessWithTimeout(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	val, ok := ProcessWithTimeout(context.Background(), ch, 100*time.Millisecond)
	if !ok || val != 42 {
		t.Errorf("expected 42, got %d (ok=%v)", val, ok)
	}

	// Таймаут на пустом канале
	emptyCh := make(chan int)
	_, ok = ProcessWithTimeout(context.Background(), emptyCh, 10*time.Millisecond)
	if ok {
		t.Errorf("expected timeout on empty channel")
	}
}

func TestCompileTimeInterface(t *testing.T) {
	var r Runner = &Service{}
	if r.Run() != "service running" {
		t.Errorf("unexpected run output: %s", r.Run())
	}
}

func TestRecreateMap(t *testing.T) {
	m := make(map[int]string, 1000)
	for i := 0; i < 1000; i++ {
		m[i] = "value"
	}
	if len(m) != 1000 {
		t.Fatalf("expected len 1000, got %d", len(m))
	}

	// Очищаем ключи через delete: длина 0, но бакеты остаются
	for k := range m {
		delete(m, k)
	}
	if len(m) != 0 {
		t.Fatalf("expected len 0, got %d", len(m))
	}

	// Пересоздаем мапу, освобождая старый hmap для GC
	m = RecreateMap(m)
	if len(m) != 0 {
		t.Fatalf("expected empty new map, got len %d", len(m))
	}
}

func BenchmarkGoodStruct_Allocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = make([]GoodStruct, 1000)
	}
}

func BenchmarkBadStruct_Allocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = make([]BadStruct, 1000)
	}
}
