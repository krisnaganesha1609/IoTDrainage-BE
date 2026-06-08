package core

import (
	"time"

	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// ProcessEvent evaluates CEP rules using in-RAM device state.
//
// Architecture note:
//   The IoT device is Source of Truth for NORMAL/WASPADA/BAHAYA status.
//   This engine only handles BLOCKAGE detection — water rising over 3
//   consecutive readings while it is NOT raining — which the firmware
//   does not compute.
//
// Concurrency: the caller (sensor.service) must NOT hold DeviceStatesMu
// when calling this function, because ProcessEvent acquires it internally.
func ProcessEvent(deviceID string, waterDistance float64, rainDetected bool) (utils.AlertType, bool) {
	// GetOrCreateDeviceState acquires and releases the lock before returning.
	state := utils.GetOrCreateDeviceState(deviceID)

	// Now acquire lock for the read-modify-write on state fields.
	utils.DeviceStatesMu.Lock()
	defer utils.DeviceStatesMu.Unlock()

	// Cooldown: skip if last alert was within 2 minutes.
	if time.Since(state.LastAlertTime) < 2*time.Minute {
		return "", false
	}

	// Rolling window of the last 5 distance readings (for trend analysis).
	state.LastDistances = append(state.LastDistances, waterDistance)
	state.LastTimesSlice = append(state.LastTimesSlice, time.Now())
	if len(state.LastDistances) > 5 {
		state.LastDistances = state.LastDistances[1:]
		state.LastTimesSlice = state.LastTimesSlice[1:]
	}

	// BLOCKAGE: water distance shrinking over 3 consecutive readings without rain.
	if len(state.LastDistances) >= 3 {
		l1 := state.LastDistances[len(state.LastDistances)-3]
		l2 := state.LastDistances[len(state.LastDistances)-2]
		l3 := state.LastDistances[len(state.LastDistances)-1]

		if l1 > l2 && l2 > l3 && !rainDetected {
			return utils.Blockage, true
		}
	}

	return "", false
}
