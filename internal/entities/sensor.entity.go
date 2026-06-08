package entities

import "time"

// SensorHistory is returned by REST endpoints to mobile / dashboard consumers.
// Timestamps are stored as UTC in InfluxDB; the mobile app translates to local TZ.
type SensorHistory struct {
	Time          time.Time `json:"timestamp"`
	DeviceID      string    `json:"device_id"`
	WaterDistance float64   `json:"water_distance"`
	WaterLevelCm  float64   `json:"water_level_cm"`
	Status        string    `json:"status"`          // NORMAL | WASPADA | BAHAYA (from IoT)
	RainDetected  bool      `json:"rain_detected"`
	SensorFlag    string    `json:"sensor_flag"`
	NextWakeupSec int64     `json:"next_wakeup_sec"`
}
