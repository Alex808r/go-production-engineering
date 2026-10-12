// Package main демонстрирует основы создания базового HTTP-сервера
// с использованием стандартной библиотеки Go (пакет net/http).
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// http.HandleFunc регистрирует функцию-обработчик для указанного пути (шаблона URL)
	// в глобальном маршрутизаторе по умолчанию — http.DefaultServeMux.
	//
	// Маршрут "/" в DefaultServeMux является префиксным (catch-all): он перехватывает
	// любые запросы, для которых нет более точного совпадения (например, /about, /users).
	http.HandleFunc("/", mainHandler)

	fmt.Println("HTTP-сервер запущен на http://localhost:8080")
	fmt.Println("Для остановки нажмите Ctrl+C")

	// http.ListenAndServe запускает TCP-слушатель на указанном адресе/порту
	// и начинает обрабатывать входящие HTTP-запросы в бесконечном цикле.
	//
	// Архитектура работы:
	// 1. Для каждого входящего TCP-соединения net/http запускает отдельную горутину:
	//    go c.serve(connCtx). Это обеспечивает высокую конкурентность «из коробки».
	// 2. Второй аргумент (handler) — маршрутизатор (http.Handler). Если передан nil,
	//    используется глобальный http.DefaultServeMux.
	// 3. Вызов является блокирующим. Если сервер аварийно завершится, он вернет не-nil ошибку.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// mainHandler — функция-обработчик (HTTP handler).
// Сигнатура совпадает с типом http.HandlerFunc: func(http.ResponseWriter, *http.Request).
//
// Параметры:
//   - w http.ResponseWriter: интерфейс для формирования ответа клиенту:
//     позволяет задавать заголовки (w.Header().Set(...)),
//     код ответа (w.WriteHeader(http.StatusOK)) и тело (w.Write([]byte)).
//   - r *http.Request: структура с метаданными входящего запроса:
//     метод (GET, POST), URL, заголовки, параметры строки запроса, тело (r.Body) и контекст (r.Context()).
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Если код ответа явно не установлен через w.WriteHeader(),
	// при первом вызове w.Write() Go автоматически отправляет HTTP 200 OK
	// и определяет заголовок Content-Type по содержимому (sniffing).
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<html><body><h2>Go Simple Web App</h2><p>Привет из net/http!</p></body></html>`))
}

