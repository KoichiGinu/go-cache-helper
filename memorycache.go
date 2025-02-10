package memorycache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	MemoryCacheDefaultExpiration      = 10 * time.Second
	MemoryCacheDefaultCleanupInterval = 60 * time.Second
)

var (
	ErrMemoryCacheDataNotFound    = fmt.Errorf("the cache data was not found.")
	ErrMemoryCacheAssertionFailed = fmt.Errorf("the cache data assertion is invalid.")
)

var goCacheClient = cache.New(MemoryCacheDefaultExpiration, MemoryCacheDefaultCleanupInterval)

type MemoryCacheProvider[T any] struct {
	CacheKey  string
	ValueType T
}

type MemoryCacheProviderInterface[T any] interface {
	Get() (T, error)
	Set(value T, ttl time.Duration)
	SetDefault(value T)
	Clear()
}

var _ MemoryCacheProviderInterface[any] = (*MemoryCacheProvider[any])(nil)

func NewMemoryCacheProvider[T any](cacheKey string, valueType T) *MemoryCacheProvider[T] {
	return &MemoryCacheProvider[T]{
		CacheKey:  cacheKey,
		ValueType: valueType,
	}
}

func (s MemoryCacheProvider[T]) Get() (T, error) {
	var result T
	rawResult, found := goCacheClient.Get(s.CacheKey)
	if !found {
		return result, ErrMemoryCacheDataNotFound
	}
	result, ok := rawResult.(T)
	if !ok {
		return result, ErrMemoryCacheAssertionFailed
	}

	return result, nil
}

func (s MemoryCacheProvider[T]) Set(value T, ttl time.Duration) {
	goCacheClient.Set(s.CacheKey, value, ttl)
}

func (s MemoryCacheProvider[T]) SetDefault(value T) {
	goCacheClient.SetDefault(s.CacheKey, value)
}

func (s MemoryCacheProvider[T]) Clear() {
	goCacheClient.Delete(s.CacheKey)
}