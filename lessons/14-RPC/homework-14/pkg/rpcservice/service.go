package rpcservice

import (
	"sync"

	"go-core-4/lessons/14-RPC/homework-14/pkg/messages"
)

// Service предоставляет RPC-сервис для работы с очередью сообщений.
// Удовлетворяет всем строгим требованиям стандартного пакета net/rpc:
// 1. Тип данных Service экспортируемый.
// 2. Методы Send и Messages экспортируемые.
// 3. Каждый метод имеет ровно два аргумента экспортируемого или встроенного типа.
// 4. Первый аргумент — аргумент функции (входящие данные).
// 5. Второй аргумент — указатель на результат вызова.
// 6. Метод возвращает error.
type Service struct {
	mu       sync.RWMutex
	messages []messages.Message
}

// New создает новый экземпляр RPC-сервиса.
func New() *Service {
	return &Service{
		messages: make([]messages.Message, 0),
	}
}

// Send позволяет клиенту отправить массив сообщений на сервер.
// Сообщения потокобезопасно добавляются в накопительный массив сервиса.
// В качестве ответа (*resp) возвращается количество успешно сохраненных сообщений.
func (s *Service) Send(req []messages.Message, resp *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = append(s.messages, req...)
	if resp != nil {
		*resp = len(req)
	}
	return nil
}

// Messages возвращает клиенту срез всех накопленных сообщений.
// Первый параметр — пустая структура struct{} (по требованию net/rpc о наличии двух аргументов).
func (s *Service) Messages(req struct{}, resp *[]messages.Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]messages.Message, len(s.messages))
	copy(res, s.messages)
	*resp = res
	return nil
}

// Count возвращает текущее количество сообщений в хранилище (вспомогательный метод).
func (s *Service) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.messages)
}
