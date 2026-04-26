package utils

import "time"

type DeviceState struct {
	LastDistances []float64
	LastTimes     []time.Time
	LastAlertTime time.Time
}

var DeviceStates = make(map[string]*DeviceState)

type AlertType string

const (
	HighWater AlertType = "HIGH_WATER"
	Blockage  AlertType = "BLOCKAGE"
	RainInfo  AlertType = "RAIN_INFO"
)
