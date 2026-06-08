package handlers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// ReceiveSensorFromMQTT subscribes to the sensor-data topic and processes
// each incoming message inside the MQTT callback.
//
// BUG FIX: The original code called ProcessSensorData() once immediately
// after Subscribe() — outside the callback — with an empty payload struct.
// This is now correctly placed inside the message callback so it runs on
// every received message.
func (h *Handler) ReceiveSensorFromMQTT(client *utils.MQTTClient, config *utils.MQTTConfig) {
	topic := config.TopicSensorData()
	log.Infof("[MQTT] Subscribing to sensor-data: %s", topic)

	client.Client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
		var payload requests.SensorDataRequest
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Errorf("[MQTT] Gagal unmarshal sensor-data: %v", err)
			return
		}
		log.Infof("[MQTT] sensor-data diterima dari %s | status=%s | level=%.1f cm",
			payload.DeviceID, payload.Status, payload.WaterLevelCm)

		if err := h.Service.ProcessSensorData(payload); err != nil {
			log.Errorf("[MQTT] Gagal proses sensor-data: %v", err)
		}
	})
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
