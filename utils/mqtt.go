package utils

import (
	"fmt"

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

func (config *MQTTConfig) InitMQTT() (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker) // fixed: only call AddBroker once
	opts.SetClientID("compro9.26.telyu-iot-drainage-be-subscriber")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &MQTTClient{Client: client}, nil
}
