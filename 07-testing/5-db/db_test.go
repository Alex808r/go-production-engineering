package db

import (
	"testing"
)

// TestProductCRUD демонстрирует сквозное тестирование полного жизненного цикла сущности:
// Create (создание) -> Read (чтение) -> Update (обновление) -> Delete (удаление).
//
// ИСПОЛЬЗОВАНИЕ ФЛАГА testing.Short():
// Стандартный пакет testing предоставляет функцию `testing.Short()`, которая возвращает true,
// если тесты запущены с флагом `go test -short`.
// В Go-сообществе принято использовать этот флаг для разграничения быстрых модульных тестов
// и тяжелых интеграционных сценариев (обращения к реальным БД, сетевым сервисам).
// Обычно пишут: `if testing.Short() { t.Skip("пропуск интеграционного теста в режиме -short") }`.
// В данном учебном примере проверка инвертирована: тест запускается только если явно передан флаг `-short`.
func TestProductCRUD(t *testing.T) {
	if !testing.Short() {
		t.Skip("пропуск демонстрационного теста CRUD (запуск только при флаге `go test -short`)")
	}

	var item = Product{
		Name:  "Пепелац",
		Price: 5, // цена в КЦ
	}

	// 1. CREATE: добавление нового товара
	err := NewProduct(item)
	if err != nil {
		t.Fatalf("ошибка создания продукта: %v", err)
	}

	// 2. READ: проверка наличия созданного товара в списке
	data, err := Products()
	if err != nil {
		t.Fatalf("ошибка получения списка продуктов: %v", err)
	}
	found := false
	for _, p := range data {
		if p.Name == item.Name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("не найден только что созданный продукт: %v", item)
	}

	// 3. UPDATE: обновление данных товара
	err = UpdateProduct(item)
	if err != nil {
		t.Fatalf("ошибка обновления продукта: %v", err)
	}
	data, err = Products()
	if err != nil {
		t.Fatalf("ошибка получения списка продуктов после обновления: %v", err)
	}
	found = false
	for _, p := range data {
		if p.Name == item.Name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("не найден обновленный продукт: %v", item)
	}

	// 4. DELETE: удаление товара
	err = DeleteProduct(item)
	if err != nil {
		t.Fatalf("ошибка удаления продукта: %v", err)
	}
	data, err = Products()
	if err != nil {
		t.Fatalf("ошибка получения списка продуктов после удаления: %v", err)
	}
	found = false
	for _, p := range data {
		if p.Name == item.Name {
			found = true
			break
		}
	}
	if found {
		t.Errorf("продукт не был удален: %v", item)
	}
}

