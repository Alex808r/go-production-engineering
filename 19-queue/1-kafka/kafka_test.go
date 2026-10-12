package main

import (
	"testing"
)

func TestNew_Validation(t *testing.T) {
	// 1. Пустые брокеры
	_, err := New(nil, "topic", "group")
	if err == nil {
		t.Error("expected error for nil brokers, got nil")
	}

	_, err = New([]string{""}, "topic", "group")
	if err == nil {
		t.Error("expected error for empty broker string, got nil")
	}

	// 2. Пустой топик
	_, err = New([]string{"localhost:9092"}, "", "group")
	if err == nil {
		t.Error("expected error for empty topic, got nil")
	}

	// 3. Пустой groupID
	_, err = New([]string{"localhost:9092"}, "topic", "")
	if err == nil {
		t.Error("expected error for empty groupID, got nil")
	}

	// 4. Корректные параметры
	client, err := New([]string{"localhost:9092"}, "topic", "group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Reader == nil || client.Writer == nil {
		t.Error("expected non-nil Reader and Writer")
	}

	_ = client.Close()
}
