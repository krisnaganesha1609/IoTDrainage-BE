package requests

// DeviceInfoRequest mirrors the maintenance uplink payload published by IoT
// to the device-info MQTT topic (sent once per cold boot or commissioning).
type DeviceInfoRequest struct {
	DeviceID  string `json:"device_id"`
	IPAddress string `json:"ip_address"`
	Location  string `json:"location"`
	Timestamp int64  `json:"timestamp"` // Unix UTC from device
}
