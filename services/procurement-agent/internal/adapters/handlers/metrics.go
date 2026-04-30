package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed, partitioned by method, path and status.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency distribution of HTTP requests in seconds.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served.",
		},
	)

	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total HTTP responses with status >= 500.",
		},
		[]string{"method", "path", "status"},
	)

	// ProcurementOperations is exported for use cases (search / decision / cache hits).
	ProcurementOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "procurement_operations_total",
			Help: "Procurement domain operations grouped by operation and result.",
		},
		[]string{"operation", "result"},
	)
)

// MetricsHandler exposes /metrics for Prometheus scraping.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// MetricsMiddleware records counters, latency and in-flight gauge per request.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		path := normalizedPath(r)
		status := strconv.Itoa(sw.status)
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		if sw.status >= 500 {
			httpErrorsTotal.WithLabelValues(r.Method, path, status).Inc()
		}
	})
}

func normalizedPath(r *http.Request) string {
	if p := r.Pattern; p != "" {
		return p
	}
	return r.URL.Path
}
