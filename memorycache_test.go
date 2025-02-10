package memorycache

import (
	"reflect"
	"testing"
	"time"

	"github.com/mattn/go-nulltype"
)

// ExampleStruct is used for testing the caching of structured data.
type ExampleStruct struct {
	ExampleInt    int
	ExampleString string
	ExampleTime   time.Time
	ExampleStruct nulltype.NullInt64
}

func TestMemoryCacheProvider_Get(t *testing.T) {
	// Subtest: Cache Hit
	t.Run("Cache Hit", func(t *testing.T) {
		key := "cache-key:example1"
		value := ExampleStruct{
			ExampleInt:    1,
			ExampleString: "test",
			ExampleTime:   time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
			ExampleStruct: nulltype.NullInt64Of(1),
		}
		ttl := 10 * time.Second

		// Create a new provider and store the value.
		provider := NewMemoryCacheProvider[ExampleStruct](key)
		provider.Set(value, ttl)

		// Retrieve the cached value.
		cachedValue, err := provider.Get()
		if err != nil {
			t.Fatalf("Cache Hit: unexpected error: %v", err)
		}
		// Ensure the retrieved value matches the original.
		if !reflect.DeepEqual(cachedValue, value) {
			t.Errorf("Cache Hit: got %+v, want %+v", cachedValue, value)
		}
	})

	// Subtest: Cache Miss
	t.Run("Cache Miss", func(t *testing.T) {
		key := "cache-key:example2"
		provider := NewMemoryCacheProvider[ExampleStruct](key)
		// Attempt to retrieve a value that was never set.
		_, err := provider.Get()
		if err == nil {
			t.Error("Cache Miss: expected error for missing cache data, got nil")
		} else if err != ErrDataNotFound {
			t.Errorf("Cache Miss: expected error %v, got %v", ErrDataNotFound, err)
		}
	})

	// Subtest: Cache Expiration
	t.Run("Cache Expiration", func(t *testing.T) {
		key := "cache-key:example3"
		value := ExampleStruct{
			ExampleInt:    1,
			ExampleString: "test",
			ExampleTime:   time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
			ExampleStruct: nulltype.NullInt64Of(1),
		}
		ttl := 1 * time.Second

		// Create a provider with a short TTL.
		provider := NewMemoryCacheProvider[ExampleStruct](key)
		provider.Set(value, ttl)

		// Immediately verify the value is cached.
		if _, err := provider.Get(); err != nil {
			t.Fatalf("Cache Expiration: unexpected error retrieving cache: %v", err)
		}

		// Wait for the cache entry to expire.
		time.Sleep(2 * time.Second)

		// Attempt to retrieve the expired entry.
		_, err := provider.Get()
		if err == nil {
			t.Error("Cache Expiration: expected error for expired cache data, got nil")
		} else if err != ErrDataNotFound {
			t.Errorf("Cache Expiration: expected error %v, got %v", ErrDataNotFound, err)
		}
	})
}
