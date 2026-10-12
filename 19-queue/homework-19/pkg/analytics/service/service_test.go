package service

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go-core-4/19-queue/homework-19/pkg/event"
)

func TestService_RecordLink_CalculatesStats(t *testing.T) {
	s := New()

	now := time.Now()
	s.RecordLink(event.LinkCreatedEvent{
		EventID:   "e1",
		LinkID:    "l1",
		Original:  "https://go.dev", // len 14
		ShortURL:  "http://localhost:8080/r/l1",
		CreatedAt: now,
	})

	stats := s.GetStats()
	if stats.TotalLinks != 1 {
		t.Fatalf("expected 1 link, got %d", stats.TotalLinks)
	}
	if stats.AvgOriginalLength != 14 {
		t.Errorf("expected avg 14, got %f", stats.AvgOriginalLength)
	}
	if stats.DomainCounts["go.dev"] != 1 {
		t.Errorf("expected domain go.dev count 1, got %d", stats.DomainCounts["go.dev"])
	}

	// Вторая ссылка
	s.RecordLink(event.LinkCreatedEvent{
		EventID:   "e2",
		LinkID:    "l2",
		Original:  "https://github.com/golang/go", // len 28
		ShortURL:  "http://localhost:8080/r/l2",
		CreatedAt: now.Add(time.Second),
	})

	stats2 := s.GetStats()
	if stats2.TotalLinks != 2 {
		t.Fatalf("expected 2 links, got %d", stats2.TotalLinks)
	}
	// Total = 14 + 28 = 42, Avg = 21
	if stats2.TotalOriginalLength != 42 || stats2.AvgOriginalLength != 21 {
		t.Errorf("expected total 42 and avg 21, got total %d, avg %f", stats2.TotalOriginalLength, stats2.AvgOriginalLength)
	}
	if stats2.MinOriginalLength != 14 || stats2.MaxOriginalLength != 28 {
		t.Errorf("min/max mismatch: min=%d, max=%d", stats2.MinOriginalLength, stats2.MaxOriginalLength)
	}
	if stats2.DomainCounts["github.com"] != 1 {
		t.Errorf("expected domain github.com count 1, got %d", stats2.DomainCounts["github.com"])
	}
}

func TestService_Concurrency(t *testing.T) {
	s := New()

	const count = 100
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(idx int) {
			defer wg.Done()
			url := fmt.Sprintf("https://domain%d.org/test", idx%5)
			s.RecordLink(event.LinkCreatedEvent{
				EventID:   fmt.Sprintf("evt-%d", idx),
				LinkID:    fmt.Sprintf("id-%d", idx),
				Original:  url,
				ShortURL:  fmt.Sprintf("http://localhost:8080/r/id-%d", idx),
				CreatedAt: time.Now(),
			})
			_ = s.GetStats()
		}(i)
	}

	wg.Wait()

	finalStats := s.GetStats()
	if finalStats.TotalLinks != count {
		t.Fatalf("expected %d links, got %d", count, finalStats.TotalLinks)
	}
	var sumDomains int64
	for _, cnt := range finalStats.DomainCounts {
		sumDomains += cnt
	}
	if sumDomains != count {
		t.Errorf("domain counts sum mismatch: expected %d, got %d", count, sumDomains)
	}
}
