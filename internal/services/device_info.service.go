package services

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
)

// ProcessDeviceInfo handles the device-info MQTT topic payload.
// It upserts IP address, location, and boot timestamp into Firestore,
// and auto-registers the device if it is not yet known.
func (s *Service) ProcessDeviceInfo(req requests.DeviceInfoRequest) *fiber.Error {
	if req.DeviceID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "device_id is required")
	}

	return s.Repo.UpsertDeviceInfo(req.DeviceID, req.IPAddress, req.Location, req.Timestamp)
}
