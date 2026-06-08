package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// UploadImage handles POST /api/image (flood event photo from ESP32-CAM).
func (h *Handler) UploadImage(c fiber.Ctx) error {
	file, err := c.FormFile("image")
	deviceID := c.FormValue("device_id")
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "failed to get image from request")
	}
	if deviceID == "" {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "device_id is required")
	}

	secureURL, erro := h.Service.UploadImage(c, file, deviceID)
	if erro != nil {
		return erro
	}
	return utils.RespondWithCreated(c, "Flood image uploaded successfully", fiber.Map{"secure_url": secureURL})
}

// UploadSnapshot handles POST /api/devices/:device_id/snapshot (daily scheduled snapshot).
// device_id is taken from the URL path param; the form still carries the image file.
func (h *Handler) UploadSnapshot(c fiber.Ctx) error {
	deviceID := c.Params("device_id")
	if deviceID == "" {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "device_id is required in URL path")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return utils.RespondWithError(c, fiber.StatusBadRequest, "failed to get image from request")
	}

	secureURL, erro := h.Service.UploadImage(c, file, deviceID)
	if erro != nil {
		return erro
	}
	return utils.RespondWithCreated(c, "Daily snapshot uploaded successfully", fiber.Map{"secure_url": secureURL})
}

// GetLatestImage handles GET /api/image.
func (h *Handler) GetLatestImage(c fiber.Ctx) error {
	image, err := h.Service.GetLatestImage()
	if err != nil {
		return err
	}
	return utils.RespondWithOK(c, "Latest image retrieved successfully", image)
}
