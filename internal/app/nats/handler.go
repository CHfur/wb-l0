package natsapp

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"order-service/internal/domain/models"
	"order-service/internal/services"
)

type Handler struct {
	service *services.Service
	logger  *zap.Logger
}

func NewHandler(service *services.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) SaveOrder(msg *stan.Msg) {
	var order models.Order

	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		h.logger.Info("Wrong data format")
	}

	err = h.service.SaveOrder(order)
	if err != nil {
		h.logger.Info(err.Error())
	}
}
