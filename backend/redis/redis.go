package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hpower2/url-shortener/config"
)

type Client struct {
	*redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return &Client{rdb}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

// Cache operations for URL shortener
func (c *Client) CacheURL(ctx context.Context, shortCode, originalURL string, expiration time.Duration) error {
	err := c.Set(ctx, shortCode, originalURL, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to cache URL: %w", err)
	}
	return nil
}

func (c *Client) GetCachedURL(ctx context.Context, shortCode string) (string, error) {
	val, err := c.Get(ctx, shortCode).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("URL not found in cache")
		}
		return "", fmt.Errorf("failed to get cached URL: %w", err)
	}
	return val, nil
}

func (c *Client) DeleteCachedURL(ctx context.Context, shortCode string) error {
	err := c.Del(ctx, shortCode).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cached URL: %w", err)
	}
	return nil
}

func (c *Client) IncrementClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("clicks:%s", shortCode)
	val, err := c.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment click count: %w", err)
	}
	return val, nil
}

func (c *Client) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("clicks:%s", shortCode)
	val, err := c.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get click count: %w", err)
	}
	return val, nil
} 