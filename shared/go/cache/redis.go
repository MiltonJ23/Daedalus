package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Flush(ctx context.Context) error
	Close() error
	HealthCheck(ctx context.Context) error
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	opt.MaxRetries = 3
	opt.PoolSize = 20
	opt.ConnMaxIdleTime = 5 * time.Minute

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (rc *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisClient) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return rc.client.Del(ctx, keys...).Err()
}

func (rc *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	val, err := rc.client.Exists(ctx, key).Result()
	return val > 0, err
}

func (rc *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return rc.client.Incr(ctx, key).Result()
}

func (rc *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	return rc.client.Decr(ctx, key).Result()
}

func (rc *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rc.client.Expire(ctx, key, expiration).Err()
}

func (rc *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	return rc.client.Keys(ctx, pattern).Result()
}

func (rc *RedisClient) Flush(ctx context.Context) error {
	return rc.client.FlushDB(ctx).Err()
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

func (rc *RedisClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return rc.client.Ping(ctx).Err()
}

func CacheKeyPattern(service string, key string) string {
	return fmt.Sprintf("cache:%s:%s", service, key)
}
