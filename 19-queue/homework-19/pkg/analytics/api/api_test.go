package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-core-4/19-queue/homework-19/pkg/analytics/service"
	"go-core-4/19-queue/homework-19/pkg/event"
)

func TestAnalyticsAPI_GetStats(t *testing.T) {
	svc := service.New()
	svc.RecordLink(event.LinkCreatedEvent{
		EventID:   "evt-1",
		LinkID:    "link-1",
		Original:  "https://github.com",
		ShortURL:  "http://localhost:8080/r/link-1",
		CreatedAt: time.Now(),
	})

	api := New(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	rec := httptest.NewRecorder()
	api.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var stats service.Stats
	if err := json.Unmarshal(rec.Body.Bytes(), &stats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if stats.TotalLinks != 1 || stats.DomainCounts["github.com"] != 1 {
		t.Errorf("stats mismatch: %+v", stats)
	}
}

func TestAnalyticsAPI_HealthAndMetrics(t *testing.T) {
	svc := service.New()
	api := New(svc)

	// Healthz
	req1 := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec1 := httptest.NewRecorder()
	api.Router().ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec1.Code)
	}

	// Ready
	req2 := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec2 := httptest.NewRecorder()
	api.Router().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec2.Code)
	}

	// Metrics
	req3 := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec3 := httptest.NewRecorder()
	api.Router().ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusOK || !strings.Contains(rec3.Body.String(), "analytics_links_total") {
		t.Errorf("expected valid metrics, got %s", rec3.Body.String())
	}
}
