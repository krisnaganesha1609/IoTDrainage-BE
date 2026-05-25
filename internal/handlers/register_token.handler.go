package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (h *Handler) RegisterToken(c fiber.Ctx) error {
	var request requests.RegisterTokenRequest
	if err := c.Bind().Body(&request); err != nil {
		errorsMap := utils.MapValidationErrors(err)
		if errorsMap != nil {
			return utils.RespondWithValidationError(c, errorsMap)
		}
		return utils.RespondWithError(c, fiber.StatusBadRequest, "Bad request")
	}
	err := h.Service.RegisterToken(c.Context(), &request)
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to register token")
	}
	return utils.RespondWithOK(c, "Token registered successfully", nil)
}
