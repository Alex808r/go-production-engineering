package membot

import (
	"go-core-4/lessons/07-testing/homework-07/pkg/crawler"
)

// Service - тестовая заглушка (Mock/Stub) поискового робота для автономного модульного тестирования.
type Service struct{}

// New создает новый экземпляр mock-робота.
func New() *Service {
	return &Service{}
}

// Scan возвращает предопределенный набор тестовых страниц в памяти.
func (s *Service) Scan(url string, depth int) ([]crawler.Document, error) {
	return []crawler.Document{
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
	}, nil
}

// BatchScan - заглушка пакета.
func (s *Service) BatchScan(urls []string, depth int, workers int) (<-chan crawler.Document, <-chan error) {
	ch := make(chan crawler.Document)
	errCh := make(chan error)
	close(ch)
	close(errCh)
	return ch, errCh
}
