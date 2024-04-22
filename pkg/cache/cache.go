package cache

import (
	"sync"

	"github.com/dgraph-io/ristretto"
)

var (
	cache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	mu sync.Mutex
)

func Set(key string, value interface{}) bool {
	mu.Lock()
	defer mu.Unlock()
	return cache.Set(key, value, 1)
}

func Get(key string) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()
	return cache.Get(key)
}
