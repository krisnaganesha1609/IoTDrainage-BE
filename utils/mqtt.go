package utils

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	Broker    string
	BaseTopic string
}

func (c *MQTTConfig) TopicSensorData() string {
	return fmt.Sprintf("%s/sensor-data", c.BaseTopic)
}

func (c *MQTTConfig) TopicSensorLog() string {
	return fmt.Sprintf("%s/sensor-log", c.BaseTopic)
}

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

func (config *MQTTConfig) InitMQTT(onConnect mqtt.OnConnectHandler) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID("compro9.26.telyu-iot-drainage-be-subscriber")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetCleanSession(false)
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(5 * time.Second)

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
