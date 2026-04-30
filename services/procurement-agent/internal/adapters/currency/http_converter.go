// Package currency provides CurrencyConverter adapters (FR-PROC-21).
//
// HTTPConverter calls a public exchange-rate API (default: exchangerate.host,
// no API key required) and caches the rate for 1h to avoid hammering the
// upstream. On any failure it falls back to a configurable static rate so the
// service never blocks a procurement search on currency lookup.
package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

const defaultEndpoint = "https://api.exchangerate.host/latest?base=USD&symbols=XAF"

type HTTPConverter struct {
	endpoint   string
	httpClient *http.Client
	fallback   float64

	mu         sync.RWMutex
	cachedAt   time.Time
	cachedRate float64
}

func NewHTTPConverter(endpoint string, fallback float64, timeout time.Duration) *HTTPConverter {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	if fallback <= 0 {
		fallback = 600
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &HTTPConverter{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: timeout},
		fallback:   fallback,
	}
}

func (c *HTTPConverter) USDToXAF(ctx context.Context, amountUSD float64) (float64, error) {
	rate, err := c.rate(ctx)
	if err != nil {
		log.Printf("currency: live rate fetch failed (%v), using fallback %.2f", err, c.fallback)
		rate = c.fallback
	}
	return round2(amountUSD * rate), nil
}

func (c *HTTPConverter) rate(ctx context.Context) (float64, error) {
	c.mu.RLock()
	if c.cachedRate > 0 && time.Since(c.cachedAt) < time.Hour {
		r := c.cachedRate
		c.mu.RUnlock()
		return r, nil
	}
	c.mu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("call upstream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("upstream returned %d", resp.StatusCode)
	}

	var payload struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}
	rate, ok := payload.Rates["XAF"]
	if !ok || rate <= 0 {
		return 0, fmt.Errorf("XAF rate missing or invalid")
	}

	c.mu.Lock()
	c.cachedRate = rate
	c.cachedAt = time.Now()
	c.mu.Unlock()
	return rate, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
