package handlers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (h *Handler) ReceiveSensorFromMQTT(client *utils.MQTTClient, config *utils.MQTTConfig) error {
	var payload requests.SensorDataRequest

	client.Client.Subscribe(config.Topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Error("Failed to unmarshal sensor data: %v", err)
			return
		}
		log.Infof("Received sensor data: %v", payload)
	})

	err := h.Service.ProcessSensorData(payload)

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetSensorHistory(c fiber.Ctx) error {
	deviceID := c.Params("device_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.Service.GetSensorHistory(deviceID, startDate, endDate)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}

	return utils.RespondWithOK(c, "Sensor history fetched successfully", data)
}

func (h *Handler) GetLatestSensorData(c fiber.Ctx) error {
	deviceID := c.Params("device_id")

	data, err := h.Service.GetLatestSensorData(deviceID)
	if err != nil {
		return utils.RespondWithError(c, err.Code, err.Message)
	}
	return utils.RespondWithOK(c, "Latest sensor data fetched successfully", data)
}
