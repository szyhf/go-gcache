package gcache

import (
	"fmt"
	"testing"
	"time"
)

func loader(key string) (string, error) {
	return fmt.Sprintf("valueFor%s", key), nil
}

func testSetCache(t *testing.T, gc Cache[string, string], numbers int) {
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}
}

func testGetCache(t *testing.T, gc Cache[string, string], numbers int) {
	for i := 0; i < numbers; i++ {
		key := fmt.Sprintf("Key-%d", i)
		v, err := gc.Get(key)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expectedV, _ := loader(key)
		if v != expectedV {
			t.Errorf("Expected value is %v, not %v", expectedV, v)
		}
	}
}

func testGetIFPresent(t *testing.T, evT string) {
	cache :=
		New[string, string](8).
			EvictType(evT).
			LoaderFunc(
				func(key string) (string, error) {
					return "value", nil
				}).
			Build()

	v, err := cache.GetIFPresent("key")
	if err != KeyNotFoundError {
		t.Errorf("err should not be %v", err)
	}

	time.Sleep(2 * time.Millisecond)

	v, err = cache.GetIFPresent("key")
	if err != nil {
		t.Errorf("err should not be %v", err)
	}
	if v != "value" {
		t.Errorf("v should not be %v", v)
	}
}

func setItemsByRange(t *testing.T, c Cache[int, int], start, end int) {
	for i := start; i < end; i++ {
		if err := c.Set(i, i); err != nil {
			t.Error(err)
		}
	}
}

func keysToMap[K comparable, V any](keys []K) map[K]struct{} {
	m := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

func checkItemsByRange(t *testing.T, keys []int, m map[int]int, size, start, end int) {
	if len(keys) != size {
		t.Fatalf("%v != %v", len(keys), size)
	} else if len(m) != size {
		t.Fatalf("%v != %v", len(m), size)
	}
	km := keysToMap[int, int](keys)
	for i := start; i < end; i++ {
		if _, ok := km[i]; !ok {
			t.Errorf("keys should contain %v", i)
		}
		v, ok := m[i]
		if !ok {
			t.Errorf("m should contain %v", i)
			continue
		}
		if v != i {
			t.Errorf("%v != %v", v, i)
			continue
		}
	}
}

func testExpiredItems(t *testing.T, evT string) {
	size := 8
	cache :=
		New[int, int](size).
			Expiration(time.Millisecond).
			EvictType(evT).
			Build()

	setItemsByRange(t, cache, 0, size)
	checkItemsByRange(t, cache.Keys(true), cache.GetALL(true), cache.Len(true), 0, size)

	time.Sleep(time.Millisecond)

	checkItemsByRange(t, cache.Keys(false), cache.GetALL(false), cache.Len(false), 0, size)

	if l := cache.Len(true); l != 0 {
		t.Fatalf("GetALL should returns no items, but got length %v", l)
	}

	cache.Set(1, 1)
	m := cache.GetALL(true)
	if len(m) != 1 {
		t.Fatalf("%v != %v", len(m), 1)
	} else if l := cache.Len(true); l != 1 {
		t.Fatalf("%v != %v", l, 1)
	}
	if m[1] != 1 {
		t.Fatalf("%v != %v", m[1], 1)
	}
}

func getSimpleEvictedFunc[K comparable, V any](t *testing.T) func(k K, v V) {
	return func(key K, value V) {
		t.Logf("Key=%v Value=%v will be evicted.\n", key, value)
	}
}

func buildTestCache[K comparable, V any](t *testing.T, tp string, size int) Cache[K, V] {
	return New[K, V](size).
		EvictType(tp).
		EvictedFunc(getSimpleEvictedFunc[K, V](t)).
		Build()
}

func buildTestLoadingCache[K comparable, V any](t *testing.T, tp string, size int, loader func(k K) (V, error)) Cache[K, V] {
	return New[K, V](size).
		EvictType(tp).
		LoaderFunc(loader).
		EvictedFunc(getSimpleEvictedFunc[K, V](t)).
		Build()
}

func buildTestLoadingCacheWithExpiration(t *testing.T, tp string, size int, ep time.Duration) Cache[string, string] {
	return New[string, string](size).
		EvictType(tp).
		Expiration(ep).
		LoaderFunc(loader).
		EvictedFunc(getSimpleEvictedFunc[string, string](t)).
		Build()
}
