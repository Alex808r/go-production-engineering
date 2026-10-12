// Package metrics реализует сбор метрик Prometheus для микросервиса shortener.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics содержит счетчики и гистограммы Prometheus для shortener.
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	CacheHitsTotal  prometheus.Counter
	CacheMissTotal  prometheus.Counter
}

// New инициализирует и регистрирует метрики shortener.
func New(namespace string) *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests by method, path and status",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "Histogram of HTTP request latencies",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		CacheHitsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_hits_total",
				Help:      "Total number of successful cache hits from memcache service",
			},
		),
		CacheMissTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_misses_total",
				Help:      "Total number of cache misses resulting in PostgreSQL queries",
			},
		),
	}
}
