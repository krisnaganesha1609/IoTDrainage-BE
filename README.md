# 🚀 IoT Drainage Monitoring Backend

Backend system untuk monitoring drainase berbasis IoT menggunakan **MQTT**, **WebSocket**, dan **InfluxDB** untuk real-time data processing dan alerting system.

---

## 🧠 System Overview

Sistem ini dirancang untuk mendeteksi potensi penyumbatan saluran air (drainase) secara real-time menggunakan sensor IoT.

### 🔄 Data Flow

```
IoT Device → MQTT → Backend → InfluxDB
                              ↓
                        WebSocket → Mobile Apps
```

---

## ⚙️ Tech Stack

* **Golang (Fiber v3)** → Backend framework
* **MQTT (EMQX Public Broker)** → IoT communication
* **InfluxDB** → Time-series database
* **WebSocket** → Real-time streaming
* **Cloudinary** → Image storage (ESP32-CAM)
* **Scalar** → HTTP API Documentation
* **AsyncAPI** → MQTT & WebSocket Documentation
* **Air** → Hot reload during development

---

## 📦 Project Structure

```
├── internal/
│   ├── handlers/
│   ├── services/
│   ├── repositories/
├── routes/
├── configs/
├── utils/
├── docs/
│   ├── swagger.json
│   ├── mqtt-docs/
├── main.go
├── air.toml
```

---

## 🔌 Environment Variables (.env)
nanti file .env nya di send di grup

```env
PORT=3000

# InfluxDB
INFLUX_URL=
INFLUX_TOKEN=
INFLUX_ORG=
INFLUX_BUCKET=

# Cloudinary
CLOUDINARY_URL=

# MQTT
MQTT_BROKER=
MQTT_TOPIC=
```

---

## ▶️ Running the Project

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Run with Air (Hot Reload)

```bash
air
```

> Air akan otomatis reload server setiap ada perubahan kode.

---

## 📡 MQTT Integration

### 🔹 Broker

```
broker.emqx.io:1883 <- dummy>
```

### 🔹 Topic

```
drainage/sensor <- dummy>
```

### 🔹 Payload Format

```json
{
  "device_id": "dev01",
  "location": "Gate A Drainage",
  "water_distance": 14.2,
  "rain_detected": true,
  "rain_intensity": 0.65
}
```

---

## 🌐 API Documentation

### 📘 HTTP API (Scalar)

```
http://localhost:3000/http-docs
```

---

### 📡 MQTT & WebSocket Docs (AsyncAPI)

```
http://localhost:3000/mqtt-docs/index.html
```

---

## 🔌 WebSocket

### Endpoint

```
ws://localhost:3000/ws
```

### Message Format

#### Sensor Data

```json
{
  "type": "sensor",
  "data": {
    "device_id": "dev01",
    "water_distance": 14.2,
    "rain_detected": true,
    "rain_intensity": 0.65,
    "timestamp": 1714039200
  }
}
```

#### Alert

```json
{
  "type": "alert",
  "data": {
    "device_id": "dev01",
    "alert": "BLOCKAGE"
  }
}
```

---

## 📊 Features

* ✅ Real-time sensor monitoring
* ✅ Flood & blockage early detection
* ✅ MQTT-based IoT ingestion
* ✅ WebSocket live streaming
* ✅ Time-series data storage (InfluxDB)
* ✅ Image upload from ESP32-CAM
* ✅ Alert system

---

## 🧪 Development Notes

* Gunakan `air.toml` untuk development agar hot reload aktif
* MQTT subscriber dijalankan sebagai goroutine:

  ```go
  go route.Handler.ReceiveSensorFromMQTT(...)
  ```
* Pastikan InfluxDB dan Cloudinary sudah dikonfigurasi sebelum run

---

## ⚠️ Important Notes

* Sistem menggunakan public MQTT broker → tidak untuk production
* Untuk production:

  * Gunakan private broker
  * Tambahkan authentication & TLS

---

## 🧠 Insight

Backend ini tidak hanya berbasis request-response, tetapi juga mengimplementasikan **event-driven architecture** menggunakan MQTT dan WebSocket untuk efisiensi dan real-time processing.

---

## 🔥 Quote

> “A system is truly real-time not when it is fast, but when it reacts at the exact moment it matters.”
