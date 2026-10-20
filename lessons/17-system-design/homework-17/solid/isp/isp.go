// Package isp демонстрирует принцип разделения интерфейса (Interface Segregation Principle) в Go.
// Формулировка: клиенты не должны зависеть от методов, которые они не используют.
// В Go идиоматично создавать узкие интерфейсы (1-2 метода, как io.Reader, io.Writer)
// и объявлять их там, где они используются (на стороне потребителя).
package isp

import (
	"context"
	"errors"
)

// Item — простая сущность данных.
type Item struct {
	ID   int64
	Name string
}

// Reader — минимальный интерфейс только для чтения данных.
type Reader interface {
	Get(ctx context.Context, id int64) (Item, error)
}

// Writer — минимальный интерфейс только для создания/обновления данных.
type Writer interface {
	Save(ctx context.Context, item Item) error
}

// Deleter — минимальный интерфейс только для удаления данных.
type Deleter interface {
	Delete(ctx context.Context, id int64) error
}

// ReadWriteDeleter объединяет интерфейсы через композицию, когда клиенту действительно нужны все операции.
type ReadWriteDeleter interface {
	Reader
	Writer
	Deleter
}

// ReadOnlyService зависит исключительно от Reader — соблюдение ISP!
// Он не может случайно изменить или удалить данные и не требует от моков реализации ненужных методов.
type ReadOnlyService struct {
	reader Reader
}

func NewReadOnlyService(r Reader) *ReadOnlyService {
	return &ReadOnlyService{reader: r}
}

func (s *ReadOnlyService) GetItemName(ctx context.Context, id int64) (string, error) {
	item, err := s.reader.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return item.Name, nil
}

// MemStorage реализует все интерфейсы, но клиенты могут использовать его ограниченно через узкие проекции.
type MemStorage struct {
	items map[int64]Item
}

func NewMemStorage() *MemStorage {
	return &MemStorage{items: make(map[int64]Item)}
}

func (m *MemStorage) Get(_ context.Context, id int64) (Item, error) {
	it, ok := m.items[id]
	if !ok {
		return Item{}, errors.New("item not found")
	}
	return it, nil
}

func (m *MemStorage) Save(_ context.Context, item Item) error {
	m.items[item.ID] = item
	return nil
}

func (m *MemStorage) Delete(_ context.Context, id int64) error {
	delete(m.items, id)
	return nil
}
