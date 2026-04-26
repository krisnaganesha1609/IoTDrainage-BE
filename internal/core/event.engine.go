package core

import (
	"time"

	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func ProcessEvent(deviceID string, waterDistance float64, rainDetected bool) (utils.AlertType, bool) {

	state, exists := utils.DeviceStates[deviceID]
	if !exists {
		state = &utils.DeviceState{}
		utils.DeviceStates[deviceID] = state
	}

	// simpan history (max 5 data)
	state.LastDistances = append(state.LastDistances, waterDistance)
	state.LastTimes = append(state.LastTimes, time.Now())

	if len(state.LastDistances) > 5 {
		state.LastDistances = state.LastDistances[1:]
		state.LastTimes = state.LastTimes[1:]
	}

	// RULE 1: High Water
	if waterDistance > 20 {
		return utils.HighWater, true
	}

	// RULE 2: Blockage Detection (TREND)
	if len(state.LastDistances) >= 3 {
		l1 := state.LastDistances[len(state.LastDistances)-3]
		l2 := state.LastDistances[len(state.LastDistances)-2]
		l3 := state.LastDistances[len(state.LastDistances)-1]

		// naik terus + tidak hujan
		if l1 < l2 && l2 < l3 && !rainDetected {
			return utils.Blockage, true
		}
	}

	if time.Since(state.LastAlertTime) < 2*time.Minute {
		return "", false
	}
	return utils.RainInfo, rainDetected
}
