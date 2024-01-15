package cache

import (
	"order-service/internal/domain/models"
)

type InMemoryOrderCache struct {
	data map[string]*models.Order
}

func NewInMemoryOrderCache(data map[string]*models.Order) *InMemoryOrderCache {
	return &InMemoryOrderCache{data: data}
}

func (i *InMemoryOrderCache) Find(key string) (*models.Order, error) {
	order, ok := i.data[key]
	if !ok {
		return nil, ErrOrderNotFoundInCache
	}

	return order, nil
}

func (i *InMemoryOrderCache) Put(key string, order *models.Order) {
	i.data[key] = order
}
