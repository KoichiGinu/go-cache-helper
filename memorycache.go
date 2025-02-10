package memorycache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// DefaultExpiration is the default duration for cache item expiration.
const DefaultExpiration = 10 * time.Second

// DefaultCleanupInterval is the interval at which expired cache items are cleaned up.
const DefaultCleanupInterval = 60 * time.Second

var (
	// ErrDataNotFound is returned when the requested cache entry does not exist.
	ErrDataNotFound = fmt.Errorf("cache data not found")
	// ErrAssertionFailed is returned when a cached value cannot be type-asserted to the expected type.
	ErrAssertionFailed = fmt.Errorf("cache data type assertion failed")
)

// cacheClient is the underlying shared in-memory cache instance.
var cacheClient = cache.New(DefaultExpiration, DefaultCleanupInterval)

// MemoryCacheProvider provides a generic interface for caching values of type T.
type MemoryCacheProvider[T any] struct {
	cacheKey string
}

// NewMemoryCacheProvider creates a new MemoryCacheProvider with the specified cache key.
func NewMemoryCacheProvider[T any](cacheKey string) *MemoryCacheProvider[T] {
	return &MemoryCacheProvider[T]{cacheKey: cacheKey}
}

// Get retrieves the cached value associated with the provider's cache key.
// If the cache entry is missing or if the type assertion fails, an error is returned.
func (m *MemoryCacheProvider[T]) Get() (T, error) {
	var result T
	raw, found := cacheClient.Get(m.cacheKey)
	if !found {
		return result, ErrDataNotFound
	}
	result, ok := raw.(T)
	if !ok {
		return result, ErrAssertionFailed
	}
	return result, nil
}

// Set caches the given value with a custom time-to-live (TTL).
func (m *MemoryCacheProvider[T]) Set(value T, ttl time.Duration) {
	cacheClient.Set(m.cacheKey, value, ttl)
}

// SetDefault caches the given value using the default expiration time.
func (m *MemoryCacheProvider[T]) SetDefault(value T) {
	cacheClient.SetDefault(m.cacheKey, value)
}

// Clear removes the cached entry associated with the provider's cache key.
func (m *MemoryCacheProvider[T]) Clear() {
	cacheClient.Delete(m.cacheKey)
}
