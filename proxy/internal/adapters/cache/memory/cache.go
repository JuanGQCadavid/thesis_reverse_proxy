package memory

import (
	"sync"
	"time"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type Cache struct {
	state map[string]*domain.CacheData
	mux   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		state: make(map[string]*domain.CacheData),
		mux:   sync.RWMutex{},
	}
}

func (cache *Cache) Get(path string) (*domain.CacheData, error) {
	data, ok := cache.state[path]

	if !ok {
		return nil, nil
	}

	if data.TTL.After(time.Now()) {
		cache.mux.Lock()
		defer cache.mux.Unlock()

		delete(cache.state, path)
		return nil, nil
	}

	return data, nil
}
func (cache *Cache) Save(path string, data *domain.CacheData) error {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.state[path] = data
	return nil
}
