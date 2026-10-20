package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Использование шаблона "Репозиторий".
type BooksRepo interface {
	GetBooks(context.Context) ([]book, error)
	AddBook(context.Context, book) (int, error)
}

// Книжка.
type book struct {
	ID          int
	Title       string
	Year        int
	Public      bool
	PublisherID int
	Publisher   string
}

func main() {
	// Объект БД (*sql.DB) представляет собой потокобезопасный пул подключений к СУБД.
	// *sql.DB — долгоживущий объект приложения (синглтон).
	// Следует создавать только один экземпляр *sql.DB на каждый источник данных
	// и передавать его компонентам приложения как зависимость.
	var db *sql.DB
	var err error

	// Подключение к БД:
	// sql.Open принимает имя зарегистрированного драйвера (в данном случае "mysql") и DSN (Data Source Name).
	// ВАЖНО: sql.Open не устанавливает немедленного TCP-соединения с сервером!
	// Он лишь валидирует синтаксис DSN и инициализирует внутренние структуры пула.
	pwd := os.Getenv("mysql_password")
	db, err = sql.Open("mysql", "root:"+pwd+"@tcp(ubuntu-server.northeurope.cloudapp.azure.com:3306)/books")
	if err != nil {
		log.Fatal(err)
	}
	// Обязательно освобождаем ресурсы пула при завершении работы программы.
	defer db.Close()

	// Настройка параметров пула соединений (Connection Pool Tuning):
	// Рекомендуется настраивать в любом продакшен-сервисе:
	// - SetMaxOpenConns: максимальное количество одновременно открытых соединений (включая занятые и простаивающие).
	// - SetMaxIdleConns: количество свободных соединений, удерживаемых в пуле для повторного использования.
	// - SetConnMaxLifetime: максимальное время жизни соединения (защищает от закрытия TCP-сессий файрволами и балансировщиками).
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	// Проверка соединения с БД:
	// Метод Ping() (или PingContext) выполняет реальный запрос к СУБД,
	// проверяя доступность сервера и валидность учетных данных.
	err = db.Ping()
	if err != nil {
		log.Println("Предупреждение: удаленная БД недоступна:", err)
		return
	}

	// Получение списка книг.
	data, err := books(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", data)

	// Добавление книг одной транзакцией.
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
	err = addBooks(db, data)
	if err != nil {
		log.Fatal(err)
	}

	// Повторное получение списка книг.
	data, err = books(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", data)
}

// books возвращает список книг из БД.
func books(db *sql.DB) ([]book, error) {
	// Выполнение запроса выборки данных:
	// Метод Query (или QueryContext) берет свободное соединение из пула,
	// отправляет параметризованный SQL-запрос и возвращает курсор *sql.Rows.
	// ВАЖНО: Знак '?' (в MySQL/SQLite) или '$1' (в PostgreSQL) защищает от SQL-инъекций.
	rows, err := db.Query(`
		SELECT id, title, year
		FROM books
		WHERE id >= ? AND id < ?;`,
		0, 10_000,
	)
	if err != nil {
		return nil, err
	}
	// defer rows.Close() критически важен: он возвращает активное соединение
	// обратно в пул sql.DB. Без вызова Close соединения в пуле исчерпаются!
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

	// Обязательная проверка: rows.Err() возвращает ошибку, если итерация
	// курсора прервалась из-за сетевого сбоя или ошибки драйвера.
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return books, nil
}

// addBooks добавляет в БД массив книг одной транзакцией.
func addBooks(db *sql.DB, books []book) error {
	// Начало транзакции (db.Begin или db.BeginTx):
	// Под транзакцию резервируется выделенное сетевое соединение из пула.
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Отмена транзакции в случае непредвиденной ошибки или паники.
	// Если транзакция уже успешно зафиксирована через tx.Commit(), вызов Rollback()
	// вернет sql.ErrTxDone и не нанесет вреда (является идемпотентным).
	defer tx.Rollback()

	// Подготовка параметризованного запроса (Prepared Statement) внутри транзакции.
	// Позволяет серверу БД один раз разобрать план выполнения и многократно применять его.
	stmt, err := tx.Prepare(`INSERT INTO books(title, year) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, book := range books {
		res, err := stmt.ExecContext(context.TODO(), book.Title, book.Year)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		fmt.Println("Создана запись с ID:", id)
	}

	// Фиксация изменений (COMMIT)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
