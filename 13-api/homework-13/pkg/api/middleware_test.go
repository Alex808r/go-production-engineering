package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewares(t *testing.T) {
	// Базовый обработчик
	var capturedReqID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReqID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Оборачиваем в цепочку middleware
	chain := headersMiddleware(requestIDMiddleware(loggingMiddleware(handler)))

	// Тест 1: без входящего X-Request-ID (должен сгенерироваться автоматически)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	chain.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Code = %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %s, want application/json; charset=utf-8", ct)
	}
	reqIDHeader := rr.Header().Get("X-Request-ID")
	if reqIDHeader == "" {
		t.Error("X-Request-ID header is empty")
	}
	if capturedReqID != reqIDHeader {
		t.Errorf("capturedReqID = %s, want %s", capturedReqID, reqIDHeader)
	}

	// Тест 2: с переданным X-Request-ID (должен сохраниться)
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set("X-Request-ID", "custom-123")
	rr2 := httptest.NewRecorder()
	chain.ServeHTTP(rr2, req2)

	if val := rr2.Header().Get("X-Request-ID"); val != "custom-123" {
		t.Errorf("X-Request-ID = %s, want custom-123", val)
	}
	if capturedReqID != "custom-123" {
		t.Errorf("capturedReqID = %s, want custom-123", capturedReqID)
	}
}
