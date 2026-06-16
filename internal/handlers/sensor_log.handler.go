package handlers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// ReceiveSensorLogFromMQTT subscribes to the sensor-log topic and stores
// each log entry in InfluxDB for the 7-day audit trail.
//
// Called from inside OnConnectHandler — re-subscribes on every (re)connect.
func (h *Handler) ReceiveSensorLogFromMQTT(client mqtt.Client, config *utils.MQTTConfig) {
	topic := config.TopicSensorLog()

	token := client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
		var payload requests.SensorLogRequest
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Errorf("[MQTT] Gagal unmarshal sensor-log: %v", err)
			return
		}
		log.Infof("[MQTT] sensor-log [%s] dari %s: %s", payload.Level, payload.DeviceID, payload.Message)

		if err := h.Service.ProcessSensorLog(payload); err != nil {
			log.Errorf("[MQTT] Gagal proses sensor-log: %v", err)
		}
	})

	token.Wait()
	if token.Error() != nil {
		log.Errorf("[MQTT] Gagal subscribe ke %s: %v", topic, token.Error())
		return
	}
	log.Infof("[MQTT] Subscribed ke sensor-log: %s", topic)
}
