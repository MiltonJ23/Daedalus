package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Status string

const (
	StatusUp      Status = "UP"
	StatusDown    Status = "DOWN"
	StatusUnknown Status = "UNKNOWN"
)

type HealthCheckResult struct {
	Status            Status            `json:"status"`
	Checks            map[string]Status `json:"checks"`
	Timestamp         time.Time         `json:"timestamp"`
	UptimeSeconds     int64             `json:"uptime_seconds"`
}

type HealthChecker interface {
	Check(ctx context.Context) (Status, error)
}

func NewHealthCheckHandler(checkers map[string]HealthChecker) http.HandlerFunc {
	startTime := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		checks := make(map[string]Status)
		overallStatus := StatusUp

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		for name, checker := range checkers {
			status, _ := checker.Check(ctx)
			checks[name] = status
			if status != StatusUp {
				overallStatus = StatusDown
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if overallStatus == StatusUp {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		result := HealthCheckResult{
			Status:        overallStatus,
			Checks:        checks,
			Timestamp:     time.Now(),
			UptimeSeconds: int64(time.Since(startTime).Seconds()),
		}
		_ = json.NewEncoder(w).Encode(result)
	}
}

type SimpleChecker struct {
	checkFn func(context.Context) error
}

func NewSimpleChecker(checkFn func(context.Context) error) HealthChecker {
	return &SimpleChecker{checkFn: checkFn}
}

func (sc *SimpleChecker) Check(ctx context.Context) (Status, error) {
	if err := sc.checkFn(ctx); err != nil {
		return StatusDown, err
	}
	return StatusUp, nil
}
