package handlers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (h *Handler) ReceiveSensorFromMQTT(client mqtt.Client, config *utils.MQTTConfig) {
	topic := config.TopicSensorData()

	token := client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
		// Notify connection watchdog that the connection is alive.
		h.MQTTClient.MarkMessageReceived()

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

	token.Wait()
	if token.Error() != nil {
		log.Errorf("[MQTT] Gagal subscribe ke %s: %v", topic, token.Error())
		return
	}
	log.Infof("[MQTT] Subscribed ke sensor-data: %s", topic)
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
