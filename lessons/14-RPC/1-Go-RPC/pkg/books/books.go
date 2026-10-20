// Package books определяет модели данных и контракты запросов
// для стандартного сетевого сервиса net/rpc в Go.
package books

// Book представляет модель сущности книги.
// Все поля структуры ОБЯЗАТЕЛЬНО должны быть экспортируемыми (с заглавной буквы),
// иначе встроенный сериализатор encoding/gob проигнорирует их при маршалинге.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Request представляет параметры RPC-запроса от клиента.
// Для соответствия контракту net/rpc структура запроса должна быть экспортируемой.
type Request struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

