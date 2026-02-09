// package handlers define 3 things:
// 1. input validation
// 2. service routing
// 3. return response
package handlers

import (
	"github.com/davesaah/fatch/internal/config"
	"github.com/davesaah/fatch/internal/services"
)

type Handler struct {
	Service *services.Service
	Config  *config.Config
}

func NewHandler(service *services.Service, config *config.Config) *Handler {
	return &Handler{Service: service, Config: config}
}
