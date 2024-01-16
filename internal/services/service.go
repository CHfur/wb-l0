package services

import (
	"context"
	"go.uber.org/zap"
	"order-service/internal/cache"
	"order-service/internal/domain/models"
)

type OrderSaver interface {
	SaveOrder(
		ctx context.Context,
		orderData *models.Order,
	) (err error)
}

type OrderProvider interface {
	Orders(ctx context.Context) (map[string]*models.Order, error)
}

type OrderCache interface {
	Find(key string) (*models.Order, error)
	Put(key string, order *models.Order)
}

type Service struct {
	logger   *zap.Logger
	saver    OrderSaver
	provider OrderProvider
	cache    OrderCache
}

func NewService(logger *zap.Logger, saver OrderSaver, provider OrderProvider, cache OrderCache) *Service {
	return &Service{logger: logger, saver: saver, provider: provider, cache: cache}
}

func (s *Service) SaveOrder(data *models.Order) error {
	err := s.saver.SaveOrder(context.Background(), data)
	if err != nil {
		return err
	}

	s.cache.Put(data.OrderUid, data)

	return nil
}

func (s *Service) FindOrder(id string) (*models.Order, error) {
	order, err := s.cache.Find(id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func LoadOrdersFromDbToCache(logger *zap.Logger, provider OrderProvider) OrderCache {
	op := "services.LoadOrdersFromDbToCache"

	log := logger.With(
		zap.String("op", op),
	)

	orders, err := provider.Orders(context.Background())
	if err != nil {
		panic(err)
	}

	log.Info("Orders successfully loaded from db to cache")

	return cache.NewInMemoryOrderCache(orders)
}
