package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v3/log"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// SubscribeAllMQTT returns an mqtt.OnConnectHandler that (re)subscribes to
// every MQTT topic the backend needs. Pass this into MQTTConfig.InitMQTT
// so it fires on the INITIAL connect and on EVERY reconnect — preventing
// the silent "subscribed once, lost on reconnect" failure mode.
func (h *Handler) SubscribeAllMQTT(config *utils.MQTTConfig) mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		log.Info("[MQTT] Terhubung ke broker — memasang ulang semua subscription...")
		h.ReceiveSensorFromMQTT(client, config)
		h.ReceiveSensorLogFromMQTT(client, config)
		h.ReceiveDeviceInfoFromMQTT(client, config)
	}
}
