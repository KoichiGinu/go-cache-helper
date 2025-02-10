package memorycache

import (
	"reflect"
	"testing"
	"time"

	"github.com/mattn/go-nulltype"
)

func TestMemoryCacheProvider_Get(t *testing.T) {
	type ExampleStruct struct {
		ExampleInt    int
		ExampleString string
		ExampleTime   time.Time
		ExampleStruct nulltype.NullInt64
	}

	key := "cache-key:example1"
	value := ExampleStruct{
		ExampleInt:    1,
		ExampleString: "test",
		ExampleTime:   time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
		ExampleStruct: nulltype.NullInt64Of(1),
	}
	ttl := 10 * time.Second

	memoryCacheProvider := NewMemoryCacheProvider(key, ExampleStruct{})
	memoryCacheProvider.Set(value, ttl)

	cachedValue, err := memoryCacheProvider.Get()
	if err != nil {
		t.Errorf("case1: MemoryCacheProvider.Get() failed. err:%s", err.Error())
	}

	if !reflect.DeepEqual(cachedValue, value) {
		t.Errorf("case1: MemoryCacheProvider.Get() = %v, want %v", cachedValue, value)
	}

	key2 := "cache-key:example2"
	memoryCacheProvider2 := NewMemoryCacheProvider(key2, ExampleStruct{})
	_, err = memoryCacheProvider2.Get()
	if err == nil {
		t.Errorf("case2: The error did not return even though there was no cache.")
	} else if err != ErrMemoryCacheDataNotFound {
		t.Errorf("case2: The invalid error returned.")
	}

	key3 := "cache-key:example3"
	ttl3 := 1 * time.Second

	memoryCacheProvider3 := NewMemoryCacheProvider(key3, ExampleStruct{})
	memoryCacheProvider3.Set(value, ttl3)

	_, err = memoryCacheProvider3.Get()
	if err != nil {
		t.Errorf("case3: MemoryCacheProvider.Get() failed. err:%s", err.Error())
	}

	time.Sleep(2 * time.Second)

	_, err = memoryCacheProvider3.Get()
	if err == nil {
		t.Errorf("case3: The error did not return even though the cache was expired.")
	} else if err != ErrMemoryCacheDataNotFound {
		t.Errorf("case3: The invalid error returned.")
	}
}
