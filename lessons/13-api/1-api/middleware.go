// Package main содержит промежуточные обработчики (Middleware).
package main

import (
	"log"
	"net/http"
)

// logMiddleware логирует каждый входящий HTTP-запрос (метод, адрес клиента, URI).
// Реализует стандартный паттерн middleware: func(http.Handler) http.Handler.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RemoteAddr, r.RequestURI)

		// Передача управления следующему обработчику в цепочке
		next.ServeHTTP(w, r)
	})
}

// headersMiddleware устанавливает обязательные заголовки для REST API
// (в частности, Content-Type: application/json) до вызова конечного хэндлера.
func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}

