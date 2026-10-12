// Package service реализует бизнес-логику сбора и агрегации аналитики по ссылкам.
package service

import (
	"math"
	"net/url"
	"strings"
	"sync"
	"time"

	"go-core-4/19-queue/homework-19/pkg/event"
)

// Stats содержит агрегированные аналитические показатели по ссылкам.
type Stats struct {
	TotalLinks          int64            `json:"total_links"`
	TotalOriginalLength int64            `json:"total_original_length"`
	AvgOriginalLength   float64          `json:"avg_original_length"`
	TotalShortLength    int64            `json:"total_short_length"`
	AvgShortLength      float64          `json:"avg_short_length"`
	MinOriginalLength   int              `json:"min_original_length"`
	MaxOriginalLength   int              `json:"max_original_length"`
	DomainCounts        map[string]int64 `json:"domain_counts"`
	LastLinkID          string           `json:"last_link_id,omitempty"`
	LastOriginal        string           `json:"last_original,omitempty"`
	LastEventAt         time.Time        `json:"last_event_at,omitempty"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

// Service инкапсулирует состояние и потокобезопасные методы сбора статистики.
type Service struct {
	mu    sync.RWMutex
	stats Stats
}

// New создает новый экземпляр сервиса аналитики.
func New() *Service {
	return &Service{
		stats: Stats{
			DomainCounts: make(map[string]int64),
			UpdatedAt:    time.Now().UTC(),
		},
	}
}

// RecordLink обновляет статистику на основе поступившего события создания ссылки.
func (s *Service) RecordLink(ev event.LinkCreatedEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	origLen := len(ev.Original)
	shortLen := len(ev.ShortURL)

	s.stats.TotalLinks++
	s.stats.TotalOriginalLength += int64(origLen)
	s.stats.TotalShortLength += int64(shortLen)

	s.stats.AvgOriginalLength = roundToTwoDecimals(float64(s.stats.TotalOriginalLength) / float64(s.stats.TotalLinks))
	s.stats.AvgShortLength = roundToTwoDecimals(float64(s.stats.TotalShortLength) / float64(s.stats.TotalLinks))

	if s.stats.MinOriginalLength == 0 || origLen < s.stats.MinOriginalLength {
		s.stats.MinOriginalLength = origLen
	}
	if origLen > s.stats.MaxOriginalLength {
		s.stats.MaxOriginalLength = origLen
	}

	domain := extractDomain(ev.Original)
	if domain != "" {
		s.stats.DomainCounts[domain]++
	}

	s.stats.LastLinkID = ev.LinkID
	s.stats.LastOriginal = ev.Original
	s.stats.LastEventAt = ev.CreatedAt
	s.stats.UpdatedAt = time.Now().UTC()
}

// GetStats возвращает снимок текущей статистики.
func (s *Service) GetStats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Глубокое копирование map для предотвращения race condition снаружи
	domainsCopy := make(map[string]int64, len(s.stats.DomainCounts))
	for k, v := range s.stats.DomainCounts {
		domainsCopy[k] = v
	}

	res := s.stats
	res.DomainCounts = domainsCopy
	return res
}

// Reset сбрасывает статистику (полезно для тестов).
func (s *Service) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats = Stats{
		DomainCounts: make(map[string]int64),
		UpdatedAt:    time.Now().UTC(),
	}
}

// extractDomain извлекает имя хоста/домена из URL без порта и www-префикса.
func extractDomain(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "unknown"
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		// Попытка разобрать без схемы
		if parsed, err = url.Parse("http://" + rawURL); err != nil || parsed.Host == "" {
			return "unknown"
		}
	}

	host := parsed.Hostname()
	host = strings.TrimPrefix(host, "www.")
	return strings.ToLower(host)
}

func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}
