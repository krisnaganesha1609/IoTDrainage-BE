package handlers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// ReceiveDeviceInfoFromMQTT subscribes to the device-info topic.
// Sent once per cold boot or commissioning; stores IP and location to Firestore.
func (h *Handler) ReceiveDeviceInfoFromMQTT(client *utils.MQTTClient, config *utils.MQTTConfig) {
	topic := config.TopicDeviceInfo()
	log.Infof("[MQTT] Subscribing to device-info: %s", topic)

	client.Client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
		var payload requests.DeviceInfoRequest
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Errorf("[MQTT] Gagal unmarshal device-info: %v", err)
			return
		}
		log.Infof("[MQTT] device-info dari %s | ip=%s | location=%s",
			payload.DeviceID, payload.IPAddress, payload.Location)

		if err := h.Service.ProcessDeviceInfo(payload); err != nil {
			log.Errorf("[MQTT] Gagal proses device-info: %v", err)
		}
	})
}
