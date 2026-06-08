package requests

// SensorDataRequest mirrors the new Smart Publisher uplink payload from IoT.
// The IoT device is Source of Truth: it computes water_level_cm and status
// on-device. Backend stores and forwards without recalculating.
type SensorDataRequest struct {
	DeviceID      string  `json:"device_id"`
	WaterDistance float64 `json:"water_distance"`
	WaterLevelCm  float64 `json:"water_level_cm"`
	Status        string  `json:"status"`         // "NORMAL" | "WASPADA" | "BAHAYA"
	RainDetected  bool    `json:"rain_detected"`
	SensorFlag    string  `json:"sensor_flag"`    // "OK" | error codes
	Timestamp     int64   `json:"timestamp"`      // Unix UTC from device
	NextWakeupSec int64   `json:"next_wakeup_sec"` // used by watchdog
}
