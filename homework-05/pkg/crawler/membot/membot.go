package membot

import (
	"go-core-4/homework-05/pkg/crawler"
)

// Service - имитация службы поискового робота для тестов и автономной работы.
type Service struct{}

// New - конструктор имитации службы поискового робота.
func New() *Service {
	return &Service{}
}

// Scan возвращает заранее подготовленный набор тестовых страниц.
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
