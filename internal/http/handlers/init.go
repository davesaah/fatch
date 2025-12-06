// package handlers define 3 things:
// 1. input validation
// 2. service routing
// 3. return response
package handlers

import "gitlab.com/davesaah/fatch/internal/services"

type Handler struct {
	Service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}
