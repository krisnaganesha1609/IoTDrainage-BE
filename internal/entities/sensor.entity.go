package entities

import "time"

type SensorHistory struct {
	Time          time.Time `json:"timestamp" `
	WaterDistance float64   `json:"water_distance" `
	RainDetected  bool      `json:"rain_detected" `
	RainIntensity float64   `json:"rain_intensity" `
}
