package membot

import (
	"go-core-4/02-syntax/homework-02/pkg/crawler"
)

// Service - имитация службы поискового робота.
type Service struct{}

// New - конструктор имитации службы поискового робота.
func New() *Service {
	s := Service{}
	return &s
}

// Scan возвращает заранее подготовленный набор данных.
func (s *Service) Scan(url string, depth int) ([]crawler.Document, error) {
	data := []crawler.Document{
		{
			ID:    0,
			URL:   "https://go.dev",
			Title: "The Go Programming Language documents and tutorials",
		},
		{
			ID:    1,
			URL:   "https://golang.org",
			Title: "The Go Programming Language official documents",
		},
		{
			ID:    2,
			URL:   "https://yandex.ru",
			Title: "Яндекс - быстрый поиск в интернете",
		},
		{
			ID:    3,
			URL:   "https://google.com",
			Title: "Google Search Engine",
		},
	}

	return data, nil
}
