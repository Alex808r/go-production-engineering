package testmain

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4"
)

// Глобальная переменная для сервиса, инициализируемая один раз для всех тестов пакета в TestMain.
var testService *Service

// TestMain — специальная функция, управляющая запуском тестов во всем пакете.
//
// ОСОБЕННОСТИ ФУНКЦИИ TestMain:
// 1. Точка входа: Если в пакете объявлена функция `func TestMain(m *testing.M)`, то вместо
//    прямого запуска тестов Go вызывает TestMain.
// 2. Фаза Setup: Выполняется до запуска тестов (создание соединений с БД, чтение конфигов,
//    применение миграций или запуск Docker-контейнеров через testcontainers).
// 3. Запуск тестов: Вызов `m.Run()` запускает все тесты (`TestXxx`) текущего пакета и возвращает
//    числовой код завершения (0 — успех, non-zero — ошибка).
// 4. Фаза Teardown: Выполняется после `m.Run()` (закрытие пулов коннектов, удаление временных файлов).
// 5. Обязательный os.Exit: TestMain должен завершаться вызовом `os.Exit(code)`, чтобы передать
//    статус выполнения в терминал и CI/CD пайплайн.
// 6. Ограничение: TestMain может быть только один на пакет.
func TestMain(m *testing.M) {
	// Setup: Попытка подключиться к тестовой БД (для демонстрации)
	conn, err := pgx.Connect(context.Background(), "postgres://user:pwd@server/database")
	if err != nil {
		log.Printf("Внимание: тестовая БД недоступна (%v). Запуск демонстрационных тестов без подключения.", err)
	} else {
		// Teardown закрытия соединения через defer сработает только если os.Exit вызывается ниже
		defer conn.Close(context.Background())
	}

	// Инициализируем общий инстанс сервиса для всех тестов
	testService = &Service{
		db: conn,
	}

	// Запуск всех тестов пакета
	code := m.Run()

	// Teardown: здесь можно выполнить очистку глобальных ресурсов перед выходом

	os.Exit(code)
}

// TestService_Products проверяет возврат списка товаров сервисом.
// Демонстрирует использование reflect.DeepEqual для глубокого сравнения срезов структур.
func TestService_Products(t *testing.T) {
	got := testService.Products()
	want := []Product{
		{
			Name:  "Компьютер",
			Price: 20_000,
		},
	}

	// В Go срезы нельзя сравнивать оператором == (только со значением nil).
	// reflect.DeepEqual выполняет глубокое рекурсивное поэлементное сравнение полей.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Service.Products() = %v, want %v", got, want)
	}
}

