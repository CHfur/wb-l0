package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRouter(controller *Controller) http.Handler {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/v1/order/{id}", controller.FindOrder).Methods(http.MethodGet, http.MethodOptions, http.MethodHead)

	return r
}
