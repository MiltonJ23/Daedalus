package tracing

import (
	"context"
	"fmt"
	"time"
)

type TracingConfig struct {
	ServiceName string
	JaegerHost  string
	Environment string
}

func Init(cfg TracingConfig) (func(context.Context) error, error) {
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	if cfg.JaegerHost == "" {
		return nil, fmt.Errorf("jaeger host is required")
	}

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return nil
	}, nil
}

func Noop() func(context.Context) error {
	return func(ctx context.Context) error {
		return nil
	}
}
