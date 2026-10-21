package http

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP method, path, status, and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		// Only log slow requests or non-200s to avoid high log volume at 1000 RPS
		if rw.statusCode >= 400 || duration > 100*time.Millisecond {
			log.Printf("[%s] %s %s - %d (%v)", r.RemoteAddr, r.Method, r.URL.Path, rw.statusCode, duration)
		}
	})
}

// RecoveryMiddleware catches panics and returns 500 error gracefully.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] %v", rec)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
