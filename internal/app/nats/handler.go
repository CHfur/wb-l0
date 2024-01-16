package natsapp

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"order-service/internal/domain/models"
	"order-service/internal/services"
)

const maxGoroutineHandle = 90

type Handler struct {
	queue   chan *models.Order
	service *services.Service
	logger  *zap.Logger
}

func NewHandler(service *services.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger, queue: make(chan *models.Order)}
}

func (h *Handler) SaveOrder(msg *stan.Msg) {
	var order *models.Order

	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		h.logger.Info("Wrong data format")
		return
	}

	h.queue <- order
}

func (h *Handler) StartHandle(ctx context.Context) {
	for i := 0; i < maxGoroutineHandle; i++ {
		go h.handle(ctx)
	}
}

func (h *Handler) handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			order := <-h.queue

			err := h.service.SaveOrder(order)
			if err != nil {
				h.logger.Info(err.Error())
			}
		}
	}
}
