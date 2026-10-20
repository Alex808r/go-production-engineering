// Package testmain демонстрирует управление жизненным циклом тестов пакета через TestMain.
package testmain

import "github.com/jackc/pgx/v4"

// Service представляет сервис управления товарами, работающий с базой данных PostgreSQL.
type Service struct {
	db *pgx.Conn // подключение к БД (в тестах инициализируется в TestMain)
}

// Product представляет карточку товара.
type Product struct {
	Name  string
	Price int
}

// Products возвращает список доступных товаров.
func (s *Service) Products() []Product {
	data := []Product{
		{
			Name:  "Компьютер",
			Price: 20_000,
		},
	}
	return data
}

