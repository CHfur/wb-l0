package cache

import (
	"order-service/internal/domain/models"
	"sync"
)

type InMemoryOrderCache struct {
	mu   sync.RWMutex
	data map[string]*models.Order
}

func NewInMemoryOrderCache(data map[string]*models.Order) *InMemoryOrderCache {
	return &InMemoryOrderCache{data: data}
}

func (i *InMemoryOrderCache) Find(key string) (*models.Order, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	order, ok := i.data[key]
	if !ok {
		return nil, ErrOrderNotFoundInCache
	}

	return order, nil
}

func (i *InMemoryOrderCache) Put(key string, order *models.Order) {
	i.mu.Lock()
	i.data[key] = order
	i.mu.Unlock()
}
