package utils

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	Broker    string
	BaseTopic string
}

// TopicSensorData is the main telemetry topic.
func (c *MQTTConfig) TopicSensorData() string {
	return fmt.Sprintf("%s/sensor-data", c.BaseTopic)
}

// TopicSensorLog is the system log topic.
func (c *MQTTConfig) TopicSensorLog() string {
	return fmt.Sprintf("%s/sensor-log", c.BaseTopic)
}

// TopicDeviceInfo is the maintenance / device-info topic.
func (c *MQTTConfig) TopicDeviceInfo() string {
	return fmt.Sprintf("%s/device-info", c.BaseTopic)
}

func LoadMQTTConfig(broker, baseTopic string) *MQTTConfig {
	return &MQTTConfig{
		Broker:    broker,
		BaseTopic: baseTopic,
	}
}

type MQTTClient struct {
	Client mqtt.Client
}

// InitMQTT connects to the broker and wires onConnect so that ALL topic
// subscriptions are (re)established every single time a connection is made —
// both on the initial connect AND on every reconnect.
//
// BUG FIX (silent data loss after reconnect):
//
//	SetAutoReconnect(true) only restores the underlying TCP connection.
//	It does NOT restore subscriptions, because CleanSession defaults to
//	true in paho — the broker discards subscription state on disconnect.
//	Previously, Subscribe() was only called once at startup (in main.go),
//	so after ANY network blip / broker idle-disconnect over the following
//	days, the backend would silently stop receiving messages — reconnected
//	at the transport level, but no longer subscribed to anything, and with
//	no log output to indicate it, because no ConnectionLostHandler existed.
//
// Fix: pass the subscribe-everything logic in as onConnect; it now reruns
// on every connect event, and we log connection-lost / reconnect attempts
// so this failure mode is visible in the logs instead of silent.
func (config *MQTTConfig) InitMQTT(onConnect mqtt.OnConnectHandler) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID("compro9.26.telyu-iot-drainage-be-subscriber")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetOnConnectHandler(onConnect)
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("[MQTT] Koneksi ke broker putus: %v — mencoba auto-reconnect...", err)
	})
	opts.SetReconnectingHandler(func(_ mqtt.Client, _ *mqtt.ClientOptions) {
		log.Println("[MQTT] Reconnecting ke broker...")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &MQTTClient{Client: client}, nil
}
