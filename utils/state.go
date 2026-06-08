package utils

import (
	"sync"
	"time"
)

type AlertType string

const (
	HighWater AlertType = "HIGH_WATER"
	Blockage  AlertType = "BLOCKAGE"
	RainInfo  AlertType = "RAIN_INFO"
)

type DeviceState struct {
	LastAlertTime  time.Time   // Cooldown tracker (in-RAM)
	LastDistances  []float64   // For blockage trend detection
	LastTimesSlice []time.Time // Time tracking for trend slice
	LastTimes      map[int]time.Time
}

// DeviceStatesMu protects DeviceStates from concurrent read/write
// (MQTT goroutine writes, HTTP handlers read).
var DeviceStatesMu sync.RWMutex
var DeviceStates = make(map[string]*DeviceState)

// GetOrCreateDeviceState safely returns the state for a device,
// creating it if it does not exist yet.
func GetOrCreateDeviceState(deviceID string) *DeviceState {
	DeviceStatesMu.Lock()
	defer DeviceStatesMu.Unlock()

	state, exists := DeviceStates[deviceID]
	if !exists {
		state = &DeviceState{
			LastTimes: make(map[int]time.Time),
		}
		DeviceStates[deviceID] = state
	}
	return state
}
