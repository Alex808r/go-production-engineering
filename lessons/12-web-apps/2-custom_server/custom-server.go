// Package main демонстрирует детальную настройку структуры http.Server
// для безопасной и надежной эксплуатации в продакшене.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	const addr = ":8080"

	// =========================================================================
	// Параметры веб-сервера (структура http.Server)
	// =========================================================================
	// Внимание: http.ListenAndServe(":8080", handler) использует конфигурацию по умолчанию,
	// где таймауты равны нулю (отсутствуют). Это делает сервер уязвимым к атакам Slowloris
	// (медленная отправка заголовков или тела клиентом блокирует горутину и сокет навсегда).
	srv := &http.Server{
		Addr: addr,

		// Handler: маршрутизатор запросов. Если nil, используется http.DefaultServeMux.
		Handler: nil,

		// ReadTimeout: максимальное время на чтение всего входящего запроса целиком,
		// включая тело запроса (body). Таймер запускается в момент принятия TCP-соединения.
		ReadTimeout: 10 * time.Second,

		// ReadHeaderTimeout: максимальное время на чтение только HTTP-заголовков (Header).
		// Рекомендуется всегда указывать для защиты от Slowloris атак.
		ReadHeaderTimeout: 3 * time.Second,

		// WriteTimeout: максимальное время на запись HTTP-ответа клиенту.
		// Таймер сбрасывается при прочтении заголовков запроса.
		WriteTimeout: 20 * time.Second,

		// IdleTimeout: максимальное время ожидания следующего запроса в Keep-Alive соединении.
		// Если не задан, используется значение ReadTimeout.
		IdleTimeout: 60 * time.Second,

		// MaxHeaderBytes: лимит на размер заголовков запроса в байтах (по умолчанию 1 МБ).
		MaxHeaderBytes: 1 << 20, // 1 Мбайт
	}

	// Создаем TCP-слушатель (net.Listener) вручную.
	// "tcp4" ограничивает прослушивание IPv4 (можно указать "tcp" для IPv4 + IPv6).
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("Ошибка создания TCP-слушателя: %v", err)
	}

	// Регистрация обработчика в http.DefaultServeMux (так как srv.Handler == nil).
	http.HandleFunc("/", mainHandler)

	fmt.Printf("Сервер запущен с кастомными таймаутами на http://localhost%s\n", addr)
	fmt.Printf("ReadTimeout: %v, WriteTimeout: %v, IdleTimeout: %v\n",
		srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout)

	// srv.Serve начинает прием соединений на готовом net.Listener.
	// Метод блокирует выполнение до ошибки или закрытия сервера (srv.Close / srv.Shutdown).
	log.Fatal(srv.Serve(listener))

	// Для включения TLS (HTTPS) используется:
	// log.Fatal(srv.ServeTLS(listener, "cert.crt", "cert.key"))
}

// mainHandler обрабатывает запросы на URL "/".
func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<html><body><h2>Go Custom Web Server</h2><p>Сервер настроен с production-таймаутами.</p></body></html>`))
}

