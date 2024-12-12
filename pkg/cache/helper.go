package cache

import (
	"context"
)

// FetchOrStore is a generic helper that handles the common cache pattern
func FetchOrStore[T any](
	ctx context.Context,
	cache *Cache,
	key string,
	fetch func() (T, error),
) (T, error) {
	var result T

	// Try from cache first
	err := cache.Get(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	// Fetch from source
	result, err = fetch()
	if err != nil {
		return result, err
	}

	// Store in cache
	_ = cache.Set(ctx, key, result, DefaultCacheTTL)
	return result, nil
}
