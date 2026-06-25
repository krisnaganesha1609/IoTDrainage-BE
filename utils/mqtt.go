package utils

import (
	"fmt"
	"log"
	"sync/atomic"
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

	// lastMessageAt tracks the Unix timestamp of the last successfully
	// received MQTT message across ALL subscribed topics.
	// Updated atomically from MQTT callback goroutines; read by the
	// connection watchdog goroutine — atomic avoids a mutex for this
	// single int64.
	lastMessageAt atomic.Int64
}

// MarkMessageReceived must be called inside every MQTT message callback
// so the connection watchdog knows the connection is still alive and
// delivering messages. Call it as the very first line of each handler.
func (c *MQTTClient) MarkMessageReceived() {
	c.lastMessageAt.Store(time.Now().Unix())
}

// StartConnectionWatchdog spawns a goroutine that checks every
// checkInterval whether any MQTT message has been received within
// the last maxSilence duration. If not — meaning the connection is
// likely a zombie (TCP alive but broker stopped delivering) — it
// forces a full disconnect + reconnect cycle so that onConnect fires
// again and resubscribes to all topics.
//
// maxSilence should be set generously above the IoT device's
// next_wakeup_sec (e.g. if device wakes every 5 min, use 10-15 min)
// to avoid false positives when the device is simply sleeping.
func (c *MQTTClient) StartConnectionWatchdog(
	config *MQTTConfig,
	onConnect mqtt.OnConnectHandler,
	maxSilence time.Duration,
	checkInterval time.Duration,
) {
	// Seed with current time so we don't false-alarm on startup.
	c.lastMessageAt.Store(time.Now().Unix())

	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		log.Printf("[MQTT-WD] Connection watchdog started (maxSilence=%v, check=%v)",
			maxSilence, checkInterval)

		for range ticker.C {
			last := time.Unix(c.lastMessageAt.Load(), 0)
			silent := time.Since(last)

			if silent <= maxSilence {
				continue
			}

			log.Printf("[MQTT-WD] Tidak ada pesan MQTT selama %v (maks toleransi %v) — "+
				"kemungkinan zombie connection, paksa reconnect...", silent.Round(time.Second), maxSilence)

			// Force-disconnect. This makes paho call ConnectionLostHandler
			// internally, then auto-reconnect, which triggers onConnect
			// and resubscribes to all topics.
			c.Client.Disconnect(500) // 500ms quiesce

			// Brief pause to let paho process the disconnect before
			// we re-initiate; avoids a tight reconnect storm.
			time.Sleep(2 * time.Second)

			token := c.Client.Connect()
			token.Wait()
			if err := token.Error(); err != nil {
				log.Printf("[MQTT-WD] Reconnect gagal: %v — akan coba lagi di interval berikutnya", err)
			} else {
				log.Println("[MQTT-WD] Reconnect berhasil")
				// Reset the timer so we don't immediately re-trigger.
				c.lastMessageAt.Store(time.Now().Unix())
			}
		}
	}()
}

func (config *MQTTConfig) InitMQTT(onConnect mqtt.OnConnectHandler) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID("compro9.26.telyu-iot-drainage-be-subscriber")
	opts.SetCleanSession(false) // broker keeps subscription state across reconnects
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
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

	mqttClient := &MQTTClient{Client: client}

	// Start the connection watchdog.
	// maxSilence: 10 menit — adjust ke next_wakeup_sec device kamu + buffer.
	// checkInterval: cek tiap 2 menit.
	mqttClient.StartConnectionWatchdog(config, onConnect, 10*time.Minute, 2*time.Minute)

	return mqttClient, nil
}
