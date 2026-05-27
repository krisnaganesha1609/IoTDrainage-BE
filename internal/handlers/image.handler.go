package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (h *Handler) UploadImage(c fiber.Ctx) error {
	file, err := c.FormFile("image")
	deviceID := c.FormValue("device_id")
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "failed to get image from request")
	}
	if deviceID == "" {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "failed to get device_id from request")
	}
	secureURL, erro := h.Service.UploadImage(c, file, deviceID)
	if erro != nil {
		return erro
	}
	return utils.RespondWithCreated(c, "", fiber.Map{"secure_url": secureURL})
}

func (h *Handler) GetLatestImage(c fiber.Ctx) error {
	image, err := h.Service.GetLatestImage()
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "latest image retrieved successfully", image)
}
