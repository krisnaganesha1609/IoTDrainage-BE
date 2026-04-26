package utils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	Broker string
	Topic  string
}

func LoadMQTTConfig(broker, topic string) *MQTTConfig {
	return &MQTTConfig{
		Broker: broker,
		Topic:  topic,
	}
}

type MQTTClient struct {
	Client mqtt.Client
}

func (config *MQTTConfig) InitMQTT() (*MQTTClient, error) {
	opts := mqtt.NewClientOptions().AddBroker(config.Broker)
	opts.AddBroker(config.Broker)
	opts.SetClientID("compro9.26.telyu-iot-drainage-be-subscriber")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &MQTTClient{Client: client}, nil
}
