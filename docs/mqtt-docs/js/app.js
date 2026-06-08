
    const schema = {
  "asyncapi": "3.1.0",
  "info": {
    "title": "Drainage Monitoring IoT System",
    "version": "1.0.0",
    "description": "Sistem monitoring drainase berbasis IoT yang menggunakan MQTT untuk ingest data sensor dan WebSocket untuk real-time streaming ke aplikasi mobile.\n\n## 🔌 MQTT Configuration (.env)\n- MQTT_BROKER: broker EMQX / private broker\n- MQTT_BASE_TOPIC: compro9.26.telyu-iot-drainage-be\n\n## ⚡ Flow\nDevice → MQTT → Backend → InfluxDB + WebSocket → Mobile App"
  },
  "servers": {
    "mqttBroker": {
      "host": "broker.emqx.io:1883",
      "protocol": "mqtt",
      "description": "Public MQTT broker (EMQX) — backend menjadi subscriber pada tiga sub-topic di bawah base topic MQTT_BASE_TOPIC."
    },
    "websocketServer": {
      "host": "localhost:3000",
      "pathname": "/ws",
      "protocol": "ws",
      "description": "WebSocket endpoint untuk client mobile/dashboard."
    }
  },
  "channels": {
    "compro9.26.telyu-iot-drainage-be/sensor-data": {
      "address": "compro9.26.telyu-iot-drainage-be/sensor-data",
      "messages": {
        "publishSensorData.message": {
          "name": "SensorDataMessage",
          "title": "Sensor Data Payload",
          "summary": "Payload data sensor dari device IoT",
          "contentType": "application/json",
          "payload": {
            "type": "object",
            "properties": {
              "device_id": {
                "type": "string",
                "example": "dev01",
                "x-parser-schema-id": "<anonymous-schema-1>"
              },
              "water_distance": {
                "type": "number",
                "example": 14.2,
                "x-parser-schema-id": "<anonymous-schema-2>"
              },
              "water_level_cm": {
                "type": "number",
                "example": 35.8,
                "x-parser-schema-id": "<anonymous-schema-3>"
              },
              "status": {
                "type": "string",
                "example": "NORMAL",
                "enum": [
                  "NORMAL",
                  "WASPADA",
                  "BAHAYA"
                ],
                "x-parser-schema-id": "<anonymous-schema-4>"
              },
              "rain_detected": {
                "type": "boolean",
                "example": true,
                "x-parser-schema-id": "<anonymous-schema-5>"
              },
              "sensor_flag": {
                "type": "string",
                "example": "OK",
                "x-parser-schema-id": "<anonymous-schema-6>"
              },
              "timestamp": {
                "type": "integer",
                "example": 1714039200,
                "x-parser-schema-id": "<anonymous-schema-7>"
              },
              "next_wakeup_sec": {
                "type": "integer",
                "example": 60,
                "x-parser-schema-id": "<anonymous-schema-8>"
              }
            },
            "required": [
              "device_id",
              "water_distance",
              "water_level_cm",
              "status",
              "rain_detected",
              "timestamp",
              "next_wakeup_sec"
            ],
            "x-parser-schema-id": "SensorDataRequest"
          },
          "x-parser-unique-object-id": "publishSensorData.message"
        }
      },
      "description": "Channel untuk pengiriman data sensor dari device IoT",
      "x-parser-unique-object-id": "compro9.26.telyu-iot-drainage-be/sensor-data"
    },
    "compro9.26.telyu-iot-drainage-be/sensor-log": {
      "address": "compro9.26.telyu-iot-drainage-be/sensor-log",
      "messages": {
        "publishSensorLog.message": {
          "name": "SensorLogMessage",
          "title": "Sensor Log Payload",
          "summary": "Payload log sistem dari device IoT",
          "contentType": "application/json",
          "payload": {
            "type": "object",
            "properties": {
              "device_id": {
                "type": "string",
                "example": "dev01",
                "x-parser-schema-id": "<anonymous-schema-9>"
              },
              "timestamp": {
                "type": "integer",
                "example": 1714039200,
                "x-parser-schema-id": "<anonymous-schema-10>"
              },
              "level": {
                "type": "string",
                "example": "INFO",
                "enum": [
                  "INFO",
                  "WARNING",
                  "ERROR"
                ],
                "x-parser-schema-id": "<anonymous-schema-11>"
              },
              "message": {
                "type": "string",
                "example": "Sensor boot completed",
                "x-parser-schema-id": "<anonymous-schema-12>"
              },
              "wake_reason": {
                "type": "string",
                "example": "timer",
                "x-parser-schema-id": "<anonymous-schema-13>"
              },
              "reset_reason": {
                "type": "string",
                "example": "power_on",
                "x-parser-schema-id": "<anonymous-schema-14>"
              },
              "active_time_ms": {
                "type": "integer",
                "example": 1200,
                "x-parser-schema-id": "<anonymous-schema-15>"
              },
              "wifi_rssi_dbm": {
                "type": "integer",
                "example": -61,
                "x-parser-schema-id": "<anonymous-schema-16>"
              },
              "network_failures": {
                "type": "integer",
                "example": 0,
                "x-parser-schema-id": "<anonymous-schema-17>"
              },
              "free_heap_bytes": {
                "type": "integer",
                "example": 123456,
                "x-parser-schema-id": "<anonymous-schema-18>"
              }
            },
            "required": [
              "device_id",
              "timestamp",
              "level",
              "message"
            ],
            "x-parser-schema-id": "SensorLogRequest"
          },
          "x-parser-unique-object-id": "publishSensorLog.message"
        }
      },
      "description": "Channel untuk pengiriman log sistem dari device IoT",
      "x-parser-unique-object-id": "compro9.26.telyu-iot-drainage-be/sensor-log"
    },
    "compro9.26.telyu-iot-drainage-be/device-info": {
      "address": "compro9.26.telyu-iot-drainage-be/device-info",
      "messages": {
        "publishDeviceInfo.message": {
          "name": "DeviceInfoMessage",
          "title": "Device Info Payload",
          "summary": "Payload info perangkat dari device IoT",
          "contentType": "application/json",
          "payload": {
            "type": "object",
            "properties": {
              "device_id": {
                "type": "string",
                "example": "dev01",
                "x-parser-schema-id": "<anonymous-schema-19>"
              },
              "ip_address": {
                "type": "string",
                "example": "192.168.1.50",
                "x-parser-schema-id": "<anonymous-schema-20>"
              },
              "location": {
                "type": "string",
                "example": "Gate A Drainage",
                "x-parser-schema-id": "<anonymous-schema-21>"
              },
              "timestamp": {
                "type": "integer",
                "example": 1714039200,
                "x-parser-schema-id": "<anonymous-schema-22>"
              }
            },
            "required": [
              "device_id",
              "ip_address",
              "location",
              "timestamp"
            ],
            "x-parser-schema-id": "DeviceInfoRequest"
          },
          "x-parser-unique-object-id": "publishDeviceInfo.message"
        }
      },
      "description": "Channel untuk pengiriman info perangkat dan commissioning",
      "x-parser-unique-object-id": "compro9.26.telyu-iot-drainage-be/device-info"
    },
    "ws": {
      "address": "ws",
      "messages": {
        "receiveRealtimeMessages.message.0": {
          "name": "WebSocketSensorMessage",
          "title": "WebSocket Sensor Message",
          "summary": "Real-time sensor data",
          "contentType": "application/json",
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "example": "sensor",
                "x-parser-schema-id": "<anonymous-schema-24>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "payload": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-data.messages.publishSensorData.message.payload",
                  "timestamp": {
                    "type": "integer",
                    "example": 1714039200,
                    "x-parser-schema-id": "<anonymous-schema-25>"
                  }
                },
                "required": [
                  "payload",
                  "timestamp"
                ],
                "x-parser-schema-id": "WebSocketSensorData"
              }
            },
            "required": [
              "type",
              "data"
            ],
            "x-parser-schema-id": "<anonymous-schema-23>"
          },
          "x-parser-unique-object-id": "receiveRealtimeMessages.message.0"
        },
        "receiveRealtimeMessages.message.1": {
          "name": "WebSocketAlertMessage",
          "title": "WebSocket Alert Message",
          "summary": "Alert ketika terjadi kondisi abnormal",
          "contentType": "application/json",
          "payload": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "example": "alert",
                "x-parser-schema-id": "<anonymous-schema-27>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "alert": {
                    "type": "string",
                    "example": "BAHAYA",
                    "x-parser-schema-id": "<anonymous-schema-28>"
                  },
                  "payload": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-data.messages.publishSensorData.message.payload",
                  "timestamp": {
                    "type": "integer",
                    "example": 1714039200,
                    "x-parser-schema-id": "<anonymous-schema-29>"
                  }
                },
                "required": [
                  "alert",
                  "payload",
                  "timestamp"
                ],
                "x-parser-schema-id": "WebSocketAlertData"
              }
            },
            "required": [
              "type",
              "data"
            ],
            "x-parser-schema-id": "<anonymous-schema-26>"
          },
          "x-parser-unique-object-id": "receiveRealtimeMessages.message.1"
        }
      },
      "description": "Channel untuk streaming data sensor dan alert ke client mobile secara real-time",
      "x-parser-unique-object-id": "ws"
    }
  },
  "operations": {
    "publishSensorData": {
      "action": "receive",
      "channel": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-data",
      "summary": "IoT device mengirim data sensor",
      "messages": [
        "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-data.messages.publishSensorData.message"
      ],
      "x-parser-unique-object-id": "publishSensorData"
    },
    "publishSensorLog": {
      "action": "receive",
      "channel": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-log",
      "summary": "IoT device mengirim sensor log",
      "messages": [
        "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-log.messages.publishSensorLog.message"
      ],
      "x-parser-unique-object-id": "publishSensorLog"
    },
    "publishDeviceInfo": {
      "action": "receive",
      "channel": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/device-info",
      "summary": "IoT device mengirim device-info",
      "messages": [
        "$ref:$.channels.compro9.26.telyu-iot-drainage-be/device-info.messages.publishDeviceInfo.message"
      ],
      "x-parser-unique-object-id": "publishDeviceInfo"
    },
    "receiveRealtimeMessages": {
      "action": "send",
      "channel": "$ref:$.channels.ws",
      "summary": "Client menerima data sensor realtime dan alert",
      "messages": [
        "$ref:$.channels.ws.messages.receiveRealtimeMessages.message.0",
        "$ref:$.channels.ws.messages.receiveRealtimeMessages.message.1"
      ],
      "x-parser-unique-object-id": "receiveRealtimeMessages"
    }
  },
  "components": {
    "messages": {
      "SensorDataMessage": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-data.messages.publishSensorData.message",
      "SensorLogMessage": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-log.messages.publishSensorLog.message",
      "DeviceInfoMessage": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/device-info.messages.publishDeviceInfo.message",
      "WebSocketSensorMessage": "$ref:$.channels.ws.messages.receiveRealtimeMessages.message.0",
      "WebSocketAlertMessage": "$ref:$.channels.ws.messages.receiveRealtimeMessages.message.1"
    },
    "schemas": {
      "SensorDataRequest": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-data.messages.publishSensorData.message.payload",
      "SensorLogRequest": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/sensor-log.messages.publishSensorLog.message.payload",
      "DeviceInfoRequest": "$ref:$.channels.compro9.26.telyu-iot-drainage-be/device-info.messages.publishDeviceInfo.message.payload",
      "WebSocketSensorData": "$ref:$.channels.ws.messages.receiveRealtimeMessages.message.0.payload.properties.data",
      "WebSocketAlertData": "$ref:$.channels.ws.messages.receiveRealtimeMessages.message.1.payload.properties.data"
    }
  },
  "x-parser-spec-parsed": true,
  "x-parser-api-version": 3,
  "x-parser-spec-stringified": true
};
    const config = {"show":{"sidebar":true},"sidebar":{"showOperations":"byDefault"}};
    const appRoot = document.getElementById('root');
    AsyncApiStandalone.render(
        { schema, config, }, appRoot
    );
  