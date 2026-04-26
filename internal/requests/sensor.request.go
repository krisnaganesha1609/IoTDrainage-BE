package requests

type SensorDataRequest struct {
	DeviceID      string  `json:"device_id"`
	Location      string  `json:"location"`
	WaterDistance float64 `json:"water_distance"`
	RainDetected  bool    `json:"rain_detected"`
	RainIntensity float64 `json:"rain_intensity"`
}
