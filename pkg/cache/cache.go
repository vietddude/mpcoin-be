package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"mpc/internal/db/redis"
	"time"
)

var DefaultCacheTTL = 24 * time.Hour

type Cache struct {
	client *redis.Client
	prefix string
}

func NewCache(client *redis.Client, prefix string) *Cache {
	return &Cache{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves data from cache and unmarshals it into the provided interface
func (c *Cache) Get(ctx context.Context, key string, value interface{}) error {
	cacheKey := c.buildKey(key)
	data, err := c.client.Get(ctx, cacheKey).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), value)
}

// Set marshals data and stores it in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	cacheKey := c.buildKey(key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, cacheKey, data, ttl).Err()
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	cacheKey := c.buildKey(key)
	return c.client.Del(ctx, cacheKey).Err()
}

// buildKey creates a namespaced cache key
func (c *Cache) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}
