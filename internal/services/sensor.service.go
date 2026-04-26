package services

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/core"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/entities"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (s *Service) ProcessSensorData(request requests.SensorDataRequest) *fiber.Error {
	if err := s.Repo.InsertSensorData(request.DeviceID, request.Location, request.WaterDistance, request.RainDetected, request.RainIntensity); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	alertType, triggered := core.ProcessEvent(request.DeviceID, request.WaterDistance, request.RainDetected)

	if !triggered {
		forwarding := entities.WebsocketMessage{
			Type: "sensor",
			Data: map[string]interface{}{
				"payload":   request,
				"timestamp": time.Now().Unix(),
			},
		}

		if time.Since(utils.DeviceStates[request.DeviceID].LastTimes[4]) > 20*time.Second {
			s.Broadcast(forwarding)
		}
	}

	if triggered {
		if err := s.Repo.InsertAlert(request.DeviceID, string(alertType)); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		msg := entities.WebsocketMessage{
			Type: "alert",
			Data: map[string]interface{}{
				"alert":     string(alertType),
				"payload":   request,
				"timestamp": time.Now().Unix(),
			},
		}
		s.Broadcast(msg)

		//TODO: Implement push notification before updating the last alert time to prevent alert flooding
		//here

		// Update the last alert time to prevent alert flooding
		utils.DeviceStates[request.DeviceID].LastAlertTime = time.Now()
	}
	return nil
}

func (s *Service) GetSensorHistory(deviceID, startDate, endDate string) ([]entities.SensorHistory, *fiber.Error) {
	result, err := s.Repo.GetSensorHistory(deviceID, startDate, endDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var sensorData []entities.SensorHistory
	for result.Next() {
		record := result.Record()
		sensorHistory := entities.SensorHistory{
			Time:          record.Time(),
			WaterDistance: record.ValueByKey("water_distance").(float64),
			RainDetected:  record.ValueByKey("rain_detected").(bool),
			RainIntensity: record.ValueByKey("rain_intensity").(float64),
		}
		sensorData = append(sensorData, sensorHistory)
	}
	return sensorData, nil
}

func (s *Service) GetLatestSensorData(deviceID string) ([]entities.SensorHistory, *fiber.Error) {
	result, err := s.Repo.GetLatestSensorData(deviceID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var sensorData []entities.SensorHistory
	for result.Next() {
		record := result.Record()
		sensorHistory := entities.SensorHistory{
			Time:          record.Time(),
			WaterDistance: record.ValueByKey("water_distance").(float64),
			RainDetected:  record.ValueByKey("rain_detected").(bool),
			RainIntensity: record.ValueByKey("rain_intensity").(float64),
		}
		sensorData = append(sensorData, sensorHistory)
	}
	return sensorData, nil
}
