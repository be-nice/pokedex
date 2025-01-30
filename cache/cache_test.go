package cache

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Stop()

	key := "test_key"
	val := []byte("test_value")

	cache.Add(key, val)

	retrieved, found := cache.Get(key)
	if !found {
		t.Fatal("expected to find the key in cache")
	}

	if string(retrieved) != string(val) {
		t.Errorf("expected %s, got %s", val, retrieved)
	}
}

func TestCache_ReapEntries(t *testing.T) {
	cache := NewCache(1 * time.Second)
	defer cache.Stop()

	key := "test_key"
	val := []byte("test_value")

	cache.Add(key, val)

	time.Sleep(2 * time.Second)

	_, found := cache.Get(key)
	if found {
		t.Fatal("expected entry to be reaped from the cache")
	}
}
