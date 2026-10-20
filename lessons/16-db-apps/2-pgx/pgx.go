package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// book представляет модель книги в реляционной СУБД.
type book struct {
	ID          int
	Title       string
	Year        int
	Public      bool
	PublisherID int
	Publisher   string
}

func main() {
	ctx := context.Background()

	// 1. Инициализация пула соединений pgxpool.Pool.
	// pgx — специализированный драйвер PostgreSQL с поддержкой бинарного протокола,
	// типов данных Postgres (массивы, jsonb, hstore) и встроенного пула соединений pgxpool.
	pwd := os.Getenv("pg_password")
	connStr := "postgres://postgres:" + pwd + "@ubuntu-server.northeurope.cloudapp.azure.com/books"

	// pgxpool.New создает пул соединений.
	// Пул управляет открытыми TCP-соединениями с сервером СУБД,
	// повторно используя их между горутинами.
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// Обязательно освобождаем ресурсы пула при завершении работы программы.
	defer db.Close()

	// 2. Получение списка книг через параметризованный запрос.
	data, err := books(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", data)

	// 3. Пакетная вставка нескольких записей в рамках единой транзакции.
	data = []book{
		{
			Title: "The Chronicles of Amber",
			Year:  1970,
		},
		{
			Title: "Dune",
			Year:  1965,
		},
	}
	err = addBooks(ctx, db, data)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Повторное получение обновленного списка книг.
	data, err = books(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", data)
}

// books возвращает список книг из БД PostgreSQL.
// Функция демонстрирует:
// - передачу контекста для контроля таймаутов;
// - параметризацию ($1, $2) для защиты от SQL-инъекций;
// - сканирование курсора rows.Scan и проверку rows.Err().
func books(ctx context.Context, db *pgxpool.Pool) ([]book, error) {
	// Query выполняет:
	// 1) захват свободного соединения из пула;
	// 2) подготовку и отправку запроса;
	// 3) возврат структуры pgx.Rows;
	// 4) соединение удерживается до вызова rows.Close().
	rows, err := db.Query(ctx, `
		SELECT id, title, year
		FROM books
		WHERE id >= $1 AND id < $2;`,
		0, 10_000,
	)
	if err != nil {
		return nil, err
	}
	// Обязательно закрываем rows, чтобы вернуть соединение обратно в пул!
	defer rows.Close()

	var books []book
	for rows.Next() {
		var b book
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Year,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	// Проверка на сетевые ошибки или сбои сервера при итерации по строкам
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// addBooks добавляет в БД массив книг одной транзакцией с использованием pgx.Batch.
// pgx.Batch позволяет отправить сразу несколько команд на сервер за один сетевой round-trip,
// что кардинально повышает пропускную способность (throughput) по сравнению с одиночными INSERT.
func addBooks(ctx context.Context, db *pgxpool.Pool, books []book) error {
	// Открываем транзакцию
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	// Гарантированный откат в случае паники или ошибки перед Commit
	defer tx.Rollback(ctx)

	// Формируем пакет команд (batch)
	batch := &pgx.Batch{}
	for _, book := range books {
		batch.Queue(`INSERT INTO books(title, year) VALUES ($1, $2)`, book.Title, book.Year)
	}

	// Отправляем весь пакет команд в транзакцию одним сетевым пакетом
	res := tx.SendBatch(ctx, batch)

	// Закрываем объект batch results (считывает все ответы от сервера)
	err = res.Close()
	if err != nil {
		return err
	}

	// Фиксируем транзакцию (COMMIT)
	return tx.Commit(ctx)
}
