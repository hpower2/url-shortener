package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hpower2/url-shortener/redis"
)

// cacheRepository implements CacheRepository interface
type cacheRepository struct {
	redis *redis.Client
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(redis *redis.Client) CacheRepository {
	return &cacheRepository{redis: redis}
}

// SetURL caches a URL mapping
func (r *cacheRepository) SetURL(ctx context.Context, shortCode, originalURL string, expiration time.Duration) error {
	key := fmt.Sprintf("url:%s", shortCode)
	return r.redis.Set(ctx, key, originalURL, expiration).Err()
}

// GetURL retrieves a cached URL
func (r *cacheRepository) GetURL(ctx context.Context, shortCode string) (string, error) {
	key := fmt.Sprintf("url:%s", shortCode)
	return r.redis.Get(ctx, key).Result()
}

// DeleteURL removes a cached URL
func (r *cacheRepository) DeleteURL(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("url:%s", shortCode)
	return r.redis.Del(ctx, key).Err()
}

// IncrementClickCount increments the click count in cache
func (r *cacheRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("clicks:%s", shortCode)
	return r.redis.Incr(ctx, key).Err()
}

// GetClickCount retrieves the click count from cache
func (r *cacheRepository) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("clicks:%s", shortCode)
	return r.redis.Get(ctx, key).Int64()
}

// Set stores a generic key-value pair
func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redis.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a generic value by key
func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

// Delete removes a generic key
func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (r *cacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.redis.Exists(ctx, key).Result()
	return result > 0, err
}
