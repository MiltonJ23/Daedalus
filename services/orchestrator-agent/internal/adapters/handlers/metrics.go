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
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orchestrator_http_requests_total",
		Help: "Total HTTP requests by method, path and status.",
	}, []string{"method", "path", "status"})

	httpLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "orchestrator_http_request_duration_seconds",
		Help:    "HTTP request latency.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	OrchestratorOps = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orchestrator_operations_total",
		Help: "Orchestrator operations (submit_goal, dispatch, …) by outcome.",
	}, []string{"operation", "outcome"})
)

func MetricsHandler() http.Handler { return promhttp.Handler() }

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(c int) { s.status = c; s.ResponseWriter.WriteHeader(c) }

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		start := time.Now()
		next.ServeHTTP(rec, r)
		httpLatency.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		httpRequests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rec.status)).Inc()
	})
}
