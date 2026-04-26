
    const schema = {
  "asyncapi": "2.6.0",
  "info": {
    "title": "Drainage Monitoring IoT System",
    "version": "1.0.0",
    "description": "Sistem monitoring drainase berbasis IoT yang menggunakan MQTT untuk ingest data sensor \ndan WebSocket untuk real-time streaming ke aplikasi mobile.\n\n## 🔌 MQTT Configuration (.env)\n- MQTT_BROKER: check-on-env\n- MQTT_TOPIC: check-on-env\n\n## ⚡ Flow\nDevice → MQTT → Backend → InfluxDB + WebSocket → Mobile App\n"
  },
  "servers": {
    "mqttBroker": {
      "url": "broker.emqx.io:1883",
      "protocol": "mqtt",
      "description": "Public MQTT broker (EMQX)"
    }
  },
  "channels": {
    "check-on-env": {
      "description": "Channel untuk pengiriman data sensor dari device IoT",
      "publish": {
        "summary": "IoT device mengirim data sensor",
        "operationId": "publishSensorData",
        "message": {
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
              "location": {
                "type": "string",
                "example": "Gate A Drainage",
                "x-parser-schema-id": "<anonymous-schema-2>"
              },
              "water_distance": {
                "type": "number",
                "example": 14.2,
                "x-parser-schema-id": "<anonymous-schema-3>"
              },
              "rain_detected": {
                "type": "boolean",
                "example": true,
                "x-parser-schema-id": "<anonymous-schema-4>"
              },
              "rain_intensity": {
                "type": "number",
                "example": 0.65,
                "x-parser-schema-id": "<anonymous-schema-5>"
              }
            },
            "required": [
              "device_id",
              "water_distance",
              "rain_detected"
            ],
            "x-parser-schema-id": "SensorDataRequest"
          }
        }
      }
    },
    "ws": {
      "description": "Channel untuk streaming data sensor ke client mobile secara real-time",
      "subscribe": {
        "summary": "Client menerima data sensor real-time dan juga alert jika terjadi kondisi abnormal",
        "operationId": "receiveSensorData",
        "message": {
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
                "x-parser-schema-id": "<anonymous-schema-7>"
              },
              "data": {
                "type": "object",
                "properties": {
                  "device_id": {
                    "type": "string",
                    "example": "dev01",
                    "x-parser-schema-id": "<anonymous-schema-8>"
                  },
                  "water_distance": {
                    "type": "number",
                    "example": 14.2,
                    "x-parser-schema-id": "<anonymous-schema-9>"
                  },
                  "rain_detected": {
                    "type": "boolean",
                    "example": true,
                    "x-parser-schema-id": "<anonymous-schema-10>"
                  },
                  "rain_intensity": {
                    "type": "number",
                    "example": 0.65,
                    "x-parser-schema-id": "<anonymous-schema-11>"
                  },
                  "timestamp": {
                    "type": "integer",
                    "example": 1714039200,
                    "x-parser-schema-id": "<anonymous-schema-12>"
                  }
                },
                "x-parser-schema-id": "SensorWSData"
              }
            },
            "x-parser-schema-id": "<anonymous-schema-6>"
          }
        }
      }
    }
  },
  "components": {
    "messages": {
      "SensorDataMessage": "$ref:$.channels.check-on-env.publish.message",
      "WebSocketSensorMessage": "$ref:$.channels.ws.subscribe.message",
      "WebSocketAlertMessage": {
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
              "x-parser-schema-id": "<anonymous-schema-14>"
            },
            "data": {
              "type": "object",
              "properties": {
                "device_id": {
                  "type": "string",
                  "example": "dev01",
                  "x-parser-schema-id": "<anonymous-schema-15>"
                },
                "alert": {
                  "type": "string",
                  "example": "BLOCKAGE",
                  "x-parser-schema-id": "<anonymous-schema-16>"
                }
              },
              "x-parser-schema-id": "AlertWSData"
            }
          },
          "x-parser-schema-id": "<anonymous-schema-13>"
        }
      }
    },
    "schemas": {
      "SensorDataRequest": "$ref:$.channels.check-on-env.publish.message.payload",
      "SensorWSData": "$ref:$.channels.ws.subscribe.message.payload.properties.data",
      "AlertWSData": "$ref:$.components.messages.WebSocketAlertMessage.payload.properties.data"
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
  