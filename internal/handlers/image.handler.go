package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (h *Handler) UploadImage(c fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to get image from request")
	}
	secureURL, erro := h.Service.UploadImage(c, file)
	if erro != nil {
		return erro
	}
	return utils.RespondWithCreated(c, "", fiber.Map{"secure_url": secureURL})
}
