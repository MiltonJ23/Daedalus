package health

import (
	"context"
	"testing"
)

func TestNewSimpleChecker(t *testing.T) {
	checker := NewSimpleChecker(func(ctx context.Context) error {
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	status, err := checker.Check(ctx)
	if err != nil {
		t.Errorf("Check failed: %v", err)
	}

	if status != StatusUp {
		t.Errorf("status = %v, want %v", status, StatusUp)
	}
}

func TestHealthCheckResult(t *testing.T) {
	checkers := map[string]HealthChecker{
		"redis": NewSimpleChecker(func(ctx context.Context) error {
			return nil
		}),
	}

	handler := NewHealthCheckHandler(checkers)
	if handler == nil {
		t.Error("handler should not be nil")
	}
}
