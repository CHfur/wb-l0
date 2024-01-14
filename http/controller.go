package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"order-service/internal/services"
)

type Controller struct {
	service *services.Service
	logger  *zap.Logger
}

func NewController(service *services.Service, logger *zap.Logger) *Controller {
	return &Controller{service: service, logger: logger}
}

func (c *Controller) FindOrder(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Received request `find order`")

	var vars = mux.Vars(r)

	order, err := c.service.FindOrder(vars["id"])

	response, _ := json.Marshal(order)
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
		response, _ = json.Marshal(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		c.logger.Info(err.Error())
	}
}
