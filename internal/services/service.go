package services

import (
	"context"
	"go.uber.org/zap"
	"order-service/internal/domain/models"
)

type OrderSaver interface {
	SaveOrder(
		ctx context.Context,
		orderData models.Order,
	) (err error)
}

type OrderProvider interface {
	Orders(ctx context.Context) (map[string]*models.Order, error)
}

type Service struct {
	logger        *zap.Logger
	orderSaver    OrderSaver
	orderProvider OrderProvider
	//TODO: extract cache
	ordersCache map[string]*models.Order
}

func NewService(logger *zap.Logger, paymentSaver OrderSaver, paymentProvider OrderProvider) *Service {
	return &Service{logger: logger, orderSaver: paymentSaver, orderProvider: paymentProvider, ordersCache: make(map[string]*models.Order)}
}

func (s *Service) SaveOrder(data models.Order) error {
	err := s.orderSaver.SaveOrder(context.Background(), data)
	if err != nil {
		return err
	}

	s.ordersCache[data.OrderUid] = &data

	return nil
}

func (s *Service) FindOrder(id string) (*models.Order, error) {
	order, ok := s.ordersCache[id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (s *Service) LoadOrdersFromDbToCache() {
	op := "services.LoadOrdersFromDbToCache"

	log := s.logger.With(
		zap.String("op", op),
	)

	orders, err := s.orderProvider.Orders(context.Background())
	if err != nil {
		panic(err)
	}

	s.ordersCache = orders

	log.Info("Orders successfully loaded from db to cache")
}
