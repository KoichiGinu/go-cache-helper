package memorycache

import (
	"time"
)

// MemoryCacheBroker provides a transparent caching layer by leveraging the underlying cache provider.
// It wraps a cache key and TTL configuration for caching operations.
type MemoryCacheBroker[T any] struct {
	key string
	ttl time.Duration
}

// NewMemoryCacheBroker creates a new CacheBroker with the specified key and TTL.
func NewMemoryCacheBroker[T any](key string, ttl time.Duration) *MemoryCacheBroker[T] {
	return &MemoryCacheBroker[T]{
		key: key,
		ttl: ttl,
	}
}

// Exec executes the provided data-fetching function through the cache.
// It first attempts to retrieve the data from the cache using the broker's key.
// If the data is not present, it calls getData to fetch the data from the source,
// then stores the result in the cache with the broker's TTL.
func (b *MemoryCacheBroker[T]) Exec(getData func() (T, error)) (T, error) {
	// Create a cache provider for the given key.
	cacheProvider := NewMemoryCacheProvider[T](b.key)

	// Attempt to retrieve data from the cache.
	if cachedData, err := cacheProvider.Get(); err == nil {
		// Data found in cache; return it.
		return cachedData, nil
	}

	// Data not found in cache; fetch it using the provided function.
	fetchedData, err := getData()
	if err != nil {
		// Return a zero value along with the error if data fetching fails.
		var zero T
		return zero, err
	}

	// Store the newly fetched data in the cache with the configured TTL.
	cacheProvider.Set(fetchedData, b.ttl)

	return fetchedData, nil
}

// ClearCache removes the cached data associated with the broker's key.
// This is mainly used for testing purposes.
func (b *MemoryCacheBroker[T]) ClearCache() {
	cacheProvider := NewMemoryCacheProvider[T](b.key)
	cacheProvider.Clear()
}
