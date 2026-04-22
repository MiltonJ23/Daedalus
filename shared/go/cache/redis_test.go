package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	if client.client == nil {
		t.Error("redis client should not be nil")
	}
}

func TestSetGet(t *testing.T) {
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := "test:key"
	value := "test:value"

	err = client.Set(ctx, key, value, 10*time.Second)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := client.Get(ctx, key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != value {
		t.Errorf("Get = %s, want %s", got, value)
	}

	client.Del(ctx, key)
}

func TestDel(t *testing.T) {
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.Set(ctx, "test:del", "value", 10*time.Second)
	err = client.Del(ctx, "test:del")
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	exists, _ := client.Exists(ctx, "test:del")
	if exists {
		t.Error("key should be deleted")
	}
}

func TestExists(t *testing.T) {
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.Set(ctx, "test:exists", "value", 10*time.Second)
	defer client.Del(ctx, "test:exists")

	exists, err := client.Exists(ctx, "test:exists")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}

	if !exists {
		t.Error("Exists = false, want true")
	}
}

func TestIncrDecr(t *testing.T) {
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := "test:counter"
	defer client.Del(ctx, key)

	val, err := client.Incr(ctx, key)
	if err != nil {
		t.Fatalf("Incr failed: %v", err)
	}

	if val != 1 {
		t.Errorf("Incr = %d, want 1", val)
	}

	val, err = client.Decr(ctx, key)
	if err != nil {
		t.Fatalf("Decr failed: %v", err)
	}

	if val != 0 {
		t.Errorf("Decr = %d, want 0", val)
	}
}

func TestHealthCheck(t *testing.T) {
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = client.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestCacheKeyPattern(t *testing.T) {
	key := CacheKeyPattern("myservice", "user:123")
	expected := "cache:myservice:user:123"
	if key != expected {
		t.Errorf("CacheKeyPattern = %s, want %s", key, expected)
	}
}
