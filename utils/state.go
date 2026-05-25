package utils

import "time"

type AlertType string

const (
	HighWater AlertType = "HIGH_WATER"
	Blockage  AlertType = "BLOCKAGE"
	RainInfo  AlertType = "RAIN_INFO"
)

type DeviceState struct {
	LastAlertTime  time.Time         // Tetap simpan di RAM untuk bypass cepat di engine
	LastDistances  []float64         // Untuk deteksi tren blockage
	LastTimesSlice []time.Time       // Untuk tracking waktu slice tren
	LastTimes      map[int]time.Time // Bawaan dari kode awal Websocket-mu
}

// Global state di RAM
var DeviceStates = make(map[string]*DeviceState)
