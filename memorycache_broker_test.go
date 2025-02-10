package memorycache

import (
	"testing"
	"time"
)

func TestMemoryCacheBroker_Exec(t *testing.T) {
	// Test that when the cache is empty, data is fetched from the origin.
	t.Run("fetches data from origin when cache is empty", func(t *testing.T) {
		cacheKey := "unique-key-for-cache-miss"
		expiration := 1 * time.Second
		broker := NewMemoryCacheBroker[any](cacheKey, expiration)

		// Flag to verify if the getData function is executed.
		called := false

		// Execute the broker; since the cache is empty, the getData function should be called.
		_, err := broker.Exec(func() (any, error) {
			called = true
			return struct{}{}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error during Exec: %v", err)
		}
		if !called {
			t.Errorf("expected getData to be called when cache is empty, but it was not")
		}
	})

	// Test that when the cache has data, the cached value is returned and the getData function is not executed.
	t.Run("returns cached data when present", func(t *testing.T) {
		cacheKey := "key-for-cached-data"
		expiration := 1000 * time.Second // Use a long TTL to keep data in cache
		broker := NewMemoryCacheBroker[any](cacheKey, expiration)

		// First call: populate the cache.
		_, err := broker.Exec(func() (any, error) {
			return struct{}{}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error during first Exec call: %v", err)
		}

		// Second call: verify that the getData function is not executed.
		called := false
		_, err = broker.Exec(func() (any, error) {
			called = true
			return struct{}{}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error during second Exec call: %v", err)
		}
		if called {
			t.Errorf("expected cached data to be returned, but getData was called")
		}
	})
}
