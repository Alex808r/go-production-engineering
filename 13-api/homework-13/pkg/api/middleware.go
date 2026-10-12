package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// headersMiddleware устанавливает заголовок Content-Type: application/json для всех ответов.
func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware генерирует сквозной Request-ID, помещает его в контекст и заголовок ответа.
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			buf := make([]byte, 8)
			_, _ = rand.Read(buf)
			reqID = fmt.Sprintf("req-%s", hex.EncodeToString(buf))
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware логирует HTTP-метод, URI, RequestID и время выполнения запроса.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		reqID := GetRequestID(r.Context())
		log.Printf("[%s] %s %s took %v", reqID, r.Method, r.URL.Path, time.Since(start))
	})
}

// GetRequestID возвращает ID запроса из контекста.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
