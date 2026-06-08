package services

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/entities"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// ProcessSensorLog stores a system-log entry received from the IoT device.
func (s *Service) ProcessSensorLog(req requests.SensorLogRequest) *fiber.Error {
	if req.DeviceID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "device_id is required")
	}

	err := s.Repo.InsertSensorLog(
		req.DeviceID,
		req.Level,
		req.Message,
		req.WakeReason,
		req.ResetReason,
		req.ActiveTimeMs,
		req.WifiRssiDbm,
		req.NetworkFailures,
		req.FreeHeapBytes,
		req.Timestamp,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// GetRecentLogs returns device logs from the last 7 days for dashboard display.
func (s *Service) GetRecentLogs(deviceID string) ([]entities.SensorLogEntity, *fiber.Error) {
	result, err := s.Repo.GetRecentLogs(deviceID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer result.Close()

	var logs []entities.SensorLogEntity
	for result.Next() {
		record := result.Record()
		logs = append(logs, entities.SensorLogEntity{
			Time:            record.Time(),
			DeviceID:        utils.SafeString(record.ValueByKey("device_id")),
			Level:           utils.SafeString(record.ValueByKey("level")),
			Message:         utils.SafeString(record.ValueByKey("message")),
			WakeReason:      utils.SafeString(record.ValueByKey("wake_reason")),
			ResetReason:     utils.SafeString(record.ValueByKey("reset_reason")),
			ActiveTimeMs:    utils.SafeInt64(record.ValueByKey("active_time_ms")),
			WifiRssiDbm:     utils.SafeInt64(record.ValueByKey("wifi_rssi_dbm")),
			NetworkFailures: utils.SafeInt64(record.ValueByKey("network_failures")),
			FreeHeapBytes:   utils.SafeInt64(record.ValueByKey("free_heap_bytes")),
		})
	}
	return logs, nil
}
