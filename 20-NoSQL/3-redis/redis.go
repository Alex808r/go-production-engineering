package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	// go-redis — официальный клиент Redis для Go с поддержкой пула соединений,
	// конвейеров (pipelines), транзакций и кластеров.
	"github.com/go-redis/redis/v8"
)

// book — сущность для кэширования в Redis.
// В отличие от реляционных таблиц, NoSQL/кэш хранит данные в денормализованном виде.
type book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	// Подключение к СУБД Redis.
	// Под капотом NewClient создает и автоматически управляет пулом TCP-соединений.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // пустая строка, если пароль не задан
		DB:       0,  // номер логической БД Redis по умолчанию (0-15)
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Проверка доступности Redis (PING/PONG)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Внимание: сервер Redis недоступен (%v). Демонстрация кода продолжается в учебных целях.", err)
		return
	}

	books := []book{
		{ID: 1, Title: "1984"},
		{ID: 2, Title: "Clean Architecture"},
	}

	// Выполнение записи в кэш с TTL
	err := setBooks(redisClient, books)
	if err != nil {
		log.Fatal(err)
	}

	// Чтение данных из кэша
	data, err := getBooks(redisClient, []int{1, 2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Книги из кэша Redis:", data)
}

// getBooks возвращает книги из кэша Redis по их идентификаторам.
// Демонстрирует паттерн работы с ключами вида "префикс:идентификатор"
// и корректную обработку промаха кэша (Cache Miss) через redis.Nil.
func getBooks(client *redis.Client, ids []int) ([]book, error) {
	var books []book
	ctx := context.Background()

	for _, id := range ids {
		key := "books:" + strconv.Itoa(id)
		val, err := client.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			// Cache Miss: ключ отсутствует в кэше или истек TTL.
			// В реальном приложении здесь запрашивается БД (PostgreSQL/MySQL),
			// после чего значение снова записывается в кэш.
			continue
		} else if err != nil {
			return nil, fmt.Errorf("ошибка чтения ключа %s из Redis: %w", key, err)
		}

		var b book
		err = json.Unmarshal([]byte(val), &b)
		if err != nil {
			return nil, fmt.Errorf("ошибка десериализации JSON для ключа %s: %w", key, err)
		}

		books = append(books, b)
	}

	return books, nil
}

// setBooks сериализует книги в JSON и сохраняет их в Redis с временем жизни (TTL).
// Установка TTL (Time-To-Live) критически важна для предотвращения исчерпания RAM
// и обеспечения автоматической инвалидации устаревших данных.
func setBooks(client *redis.Client, books []book) error {
	ctx := context.Background()

	for _, b := range books {
		// Использование двоеточия ':' — общепринятое соглашение о пространствах имен ключей в Redis
		key := "books:" + strconv.Itoa(b.ID)

		val, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("ошибка сериализации книги в JSON: %w", err)
		}

		// Запись в Redis со временем жизни 10 минут
		err = client.Set(ctx, key, string(val), 10*time.Minute).Err()
		if err != nil {
			return fmt.Errorf("ошибка записи ключа %s в Redis: %w", key, err)
		}
	}

	return nil
}
