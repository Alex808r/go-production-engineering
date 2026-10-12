// Package main демонстрирует работу с HTTP-клиентом в Go (пакет net/http):
// от базового http.Get до тонкой настройки http.Client с таймаутами,
// контекстом отмены и кастомным http.Transport для пулинга соединений.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8080"

	// =========================================================================
	// 1. Простейший вариант GET-запроса через http.Get(url)
	// =========================================================================
	// Внимание: http.Get использует http.DefaultClient, у которого НЕТ таймаута!
	// Если удаленный сервер зависнет или не закроет TCP-соединение, запрос
	// может заблокировать горутину навсегда. Для продакшена всегда задавайте таймаут.
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Сервер не доступен (проверьте, запущен ли simple-server): %v", err)
	} else {
		// КРИТИЧЕСКИ ВАЖНО: Тело ответа resp.Body обязательно должно быть закрыто!
		// Если не закрыть resp.Body:
		// 1. Утекут системные ресурсы (файловые дескрипторы сокетов).
		// 2. TCP-соединение не вернется в пул (Keep-Alive пул клиента не сможет его переиспользовать).
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("1. Результат простого http.Get:\n%s\n\n", string(body))
	}

	// =========================================================================
	// 2. Продакшен-запрос: Контекст с таймаутом + кастомный http.Client
	// =========================================================================

	// Контекст позволяет автоматически прервать HTTP-запрос при превышении лимита времени
	// или при отмене операции родительской задачей (graceful cancellation).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// http.NewRequestWithContext связывает запрос с контекстом.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// http.Transport управляет низкоуровневыми сетевыми соединениями (TCP/TLS):
	// - Пулинг соединений (Keep-Alive): переиспользование открытых TCP сокетов для снижения задержек.
	// - Ограничение одновременных соединений (MaxIdleConns, MaxIdleConnsPerHost).
	// - Таймауты рукопожатий (TLSHandshakeTimeout, DialContext Timeout).
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // Поддержка системных прокси (HTTP_PROXY, HTTPS_PROXY)
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second, // Таймаут установки TCP-соединения
			KeepAlive: 30 * time.Second, // Интервал probe-пакетов TCP Keep-Alive
		}),
		ForceAttemptHTTP2:     true,             // Автоматический апгрейд до HTTP/2 при поддержке сервером
		MaxIdleConns:          100,              // Общий лимит свободных соединений в пуле
		MaxIdleConnsPerHost:   10,               // Лимит свободных соединений на один хост
		IdleConnTimeout:       90 * time.Second, // Время жизни простаивающего соединения в пуле
		TLSHandshakeTimeout:   10 * time.Second, // Максимальное время на TLS-рукопожатие
		ExpectContinueTimeout: 1 * time.Second,  // Ожидание ответа 100-continue для POST/PUT с большими телами
	}

	// http.Client объединяет Transport и верхнеуровневый Timeout.
	// Timeout включает в себя ВСЁ: DNS-резолв, TCP connect, TLS, отправку запроса и чтение ВСЕГО тела!
	client := http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	// Выполняем запрос через клиент
	resp2, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка client.Do: %v", err)
		return
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("2. Результат запроса с кастомным клиентом:\n%s\n", string(body2))
}

// defaultTransportDialContext возвращает функцию дозвона с заданным net.Dialer.
func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}

