package requests

// SensorLogRequest mirrors the system-log uplink payload published by IoT
// to the sensor-log MQTT topic.
type SensorLogRequest struct {
	DeviceID        string `json:"device_id"`
	Timestamp       int64  `json:"timestamp"`        // Unix UTC from device
	Level           string `json:"level"`            // "INFO" | "WARNING" | "ERROR"
	Message         string `json:"message"`
	WakeReason      string `json:"wake_reason"`
	ResetReason     string `json:"reset_reason"`
	ActiveTimeMs    int64  `json:"active_time_ms"`
	WifiRssiDbm     int64  `json:"wifi_rssi_dbm"`
	NetworkFailures int64  `json:"network_failures"`
	FreeHeapBytes   int64  `json:"free_heap_bytes"`
}
