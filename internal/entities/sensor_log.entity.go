package entities

import "time"

// SensorLogEntity represents a single system-log entry stored in InfluxDB.
type SensorLogEntity struct {
	Time            time.Time `json:"timestamp"`
	DeviceID        string    `json:"device_id"`
	Level           string    `json:"level"`           // INFO | WARNING | ERROR
	Message         string    `json:"message"`
	WakeReason      string    `json:"wake_reason"`
	ResetReason     string    `json:"reset_reason"`
	ActiveTimeMs    int64     `json:"active_time_ms"`
	WifiRssiDbm     int64     `json:"wifi_rssi_dbm"`
	NetworkFailures int64     `json:"network_failures"`
	FreeHeapBytes   int64     `json:"free_heap_bytes"`
}
