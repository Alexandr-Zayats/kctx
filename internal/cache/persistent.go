package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type item struct {
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	store = map[string]item{}
	mu    sync.Mutex
	once  sync.Once
)

func cacheFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kctx", "cache.json")
}

func loadOnce() {
	once.Do(func() {
		data, err := os.ReadFile(cacheFile())
		if err != nil {
			return
		}
		_ = json.Unmarshal(data, &store)
	})
}

func save() {
	_ = os.MkdirAll(filepath.Dir(cacheFile()), 0755)
	data, _ := json.MarshalIndent(store, "", "  ")
	_ = os.WriteFile(cacheFile(), data, 0644)
}

func Get(key string) ([]byte, bool) {
	loadOnce()

	mu.Lock()
	defer mu.Unlock()

	it, ok := store[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(it.ExpiresAt) {
		delete(store, key)
		save()
		return nil, false
	}

	return it.Data, true
}

func Set(key string, data []byte, ttl time.Duration) {
	loadOnce()

	mu.Lock()
	store[key] = item{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	mu.Unlock()

	save()
}
