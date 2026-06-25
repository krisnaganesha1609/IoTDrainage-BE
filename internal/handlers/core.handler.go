package handlers

import (
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/services"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

type Handler struct {
	Service *services.Service
	// MQTTClient reference is needed so MQTT message callbacks can call
	// MarkMessageReceived() to keep the connection watchdog informed.
	MQTTClient *utils.MQTTClient
}

func InitializeHandler(service *services.Service) *Handler {
	return &Handler{
		Service: service,
		// MQTTClient is injected after MQTT init via SetMQTTClient().
		// This two-step init is necessary because the Handler must exist
		// before InitMQTT (to build the OnConnectHandler), but MQTTClient
		// only exists after InitMQTT returns.
	}
}

// SetMQTTClient injects the MQTTClient reference after MQTT initialisation.
// Call this in main() right after mqttConfig.InitMQTT() returns.
func (h *Handler) SetMQTTClient(client *utils.MQTTClient) {
	h.MQTTClient = client
}
