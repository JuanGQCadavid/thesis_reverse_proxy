package memory

import (
	"log"
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

func (cache *Cache) Get(key domain.HttpStatusLineMultipart) (*domain.CacheData, error) {
	log.Printf("Getting cache data for path %s", key.ToString())
	data, ok := cache.state[key.ToString()]

	if !ok {
		return nil, nil
	}

	if data.TTL.Before(time.Now()) {
		cache.mux.Lock()
		defer cache.mux.Unlock()

		delete(cache.state, key.ToString())
		return nil, nil
	}

	return data, nil
}
func (cache *Cache) Save(key domain.HttpStatusLineMultipart, data *domain.CacheData) error {
	log.Printf("Saving cache data for path %s", key.ToString())
	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.state[key.ToString()] = data
	return nil
}
