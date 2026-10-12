package main

import (
	"errors"
	"log"
)

// Пример иллюстрирует принцип инверсии зависимостей (Dependency Inversion Principle, DIP):
// 1. Модули верхнего уровня (Server, logMsg) не зависят от деталей реализации (DBLogger).
//    Они зависят от абстракции — интерфейса Logger.
// 2. Детали реализации (DBLogger, MemLogger, CustomLogger) зависят от абстракции.
// 3. Это обеспечивает слабую связанность (loose coupling) и делает код легко тестируемым.

// Logger задает контракт для логирования сообщений.
type Logger interface {
	Log(string) error
}

// DBLogger — конкретная реализация логера, пишущая в базу данных.
type DBLogger struct {
	connString string // строка для подключения к БД
}

// Log записывает сообщение в журнал в БД.
func (dbl *DBLogger) Log(msg string) error {
	if dbl.connString == "" {
		return errors.New("БД недоступна")
	}
	// выполнение записи в БД...
	return nil
}

// CustomLogger демонстрирует паттерн "Декоратор / Fallback":
// Он встраивает интерфейс Logger и переопределяет метод Log.
type CustomLogger struct {
	Logger // вложенный логер (может быть nil)
}

// Log реализует интерфейс Logger:
// Если передан вложенный логер (c.Logger != nil), делегирует вызов ему.
// Если вложенный логер не задан (nil), делает fallback в стандартный log.Println.
// ВАЖНО: без этой проверки вызов c.Logger.Log(msg) привел бы к панике (nil pointer dereference)!
func (c *CustomLogger) Log(msg string) error {
	if c.Logger != nil {
		return c.Logger.Log(msg)
	}
	log.Println(msg)
	return nil
}

// MemLogger — тестовая заглушка (Mock / In-memory spy) для модульных тестов.
// Позволяет тестировать вызывающий код без реальных баз данных и сетевых запросов.
type MemLogger struct{}

// Log ничего не делает.
// В Go разрешено опускать имена получателя и аргументов: func (*MemLogger) Log(string) error,
// если они не используются в теле метода (подавляет предупреждения линтера).
func (*MemLogger) Log(string) error {
	return nil
}

// Server — сервис бизнес-логики. Зависит ТОЛЬКО от интерфейса Logger, а не от БД.
type Server struct {
	log Logger
}

// New — конструктор сервера с внедрением зависимости (Dependency Injection через конструктор).
func New(log Logger) *Server {
	return &Server{log: log}
}

// logMsg — высокоуровневая функция, принимающая абстрактный интерфейс Logger.
func logMsg(l Logger, msg string) error {
	return l.Log(msg)
}

func main() {
	// Создаем логер с fallback-поведением (пишет в stdout)
	l := new(CustomLogger)
	err := logMsg(l, "Тестовое сообщение")
	if err != nil {
		log.Println(err)
		return
	}

	// Инициализируем сервер, передавая ему логер через интерфейс (DIP)
	srv := New(l)
	_ = srv
}
