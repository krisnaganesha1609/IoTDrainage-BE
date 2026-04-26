package handlers

import (
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/services"
)

type Handler struct {
	Service *services.Service
}

func InitializeHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}
