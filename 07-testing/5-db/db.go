// Package db демонстрирует прототип слоя доступа к данным (DAO / Repository) для тестирования CRUD.
package db

// Product представляет сущность товара в хранилище.
type Product struct {
	Name  string
	Price int
}

// Products считывает список всех товаров.
func Products() ([]Product, error) {
	data := []Product{
		{
			Name:  "Пепелац",
			Price: 5, // цена в КЦ
		},
	}
	return data, nil
}

// NewProduct создает новую запись товара в хранилище (Create).
func NewProduct(Product) error {
	return nil
}

// UpdateProduct обновляет данные существующего товара (Update).
func UpdateProduct(Product) error {
	return nil
}

// DeleteProduct удаляет товар из хранилища (Delete).
func DeleteProduct(Product) error {
	return nil
}

