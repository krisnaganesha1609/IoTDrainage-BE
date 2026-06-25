package handlers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func (h *Handler) ReceiveDeviceInfoFromMQTT(client mqtt.Client, config *utils.MQTTConfig) {
	topic := config.TopicDeviceInfo()

	token := client.Subscribe(topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
		h.MQTTClient.MarkMessageReceived()

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

	token.Wait()
	if token.Error() != nil {
		log.Errorf("[MQTT] Gagal subscribe ke %s: %v", topic, token.Error())
		return
	}
	log.Infof("[MQTT] Subscribed ke device-info: %s", topic)
}
