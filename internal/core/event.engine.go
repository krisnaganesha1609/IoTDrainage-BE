package core

import (
	"time"

	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

func ProcessEvent(deviceID string, waterDistance float64, rainDetected bool) (utils.AlertType, bool) {
	state, exists := utils.DeviceStates[deviceID]
	if !exists {
		state = &utils.DeviceState{
			// Pastikan map diinisialisasi agar tidak panic assignment
			LastTimes: make(map[int]time.Time),
		}
		utils.DeviceStates[deviceID] = state
	}

	// -------------------------------------------------------------------------
	// CRITICAL FIX 1: Pindahkan cooldown check ke paling atas!
	// Jika belum lewat 2 menit dari alert terakhir, jangan proses alert baru.
	// -------------------------------------------------------------------------
	if time.Since(state.LastAlertTime) < 2*time.Minute {
		return "", false
	}

	// Simpan history tren (max 5 data) tetap di RAM
	state.LastDistances = append(state.LastDistances, waterDistance)
	// Catatan: Karena tipe di kodenya slice, kita append time.Now() biasa
	state.LastTimesSlice = append(state.LastTimesSlice, time.Now())

	if len(state.LastDistances) > 5 {
		state.LastDistances = state.LastDistances[1:]
		state.LastTimesSlice = state.LastTimesSlice[1:]
	}

	// -------------------------------------------------------------------------
	// RULE 1: High Water (Asumsi: air naik = jarak ke sensor makin dekat/kecil)
	// Misal batas aman adalah jika jarak air ke sensor kurang dari 20 cm
	// -------------------------------------------------------------------------
	if waterDistance < 20 {
		return utils.HighWater, true
	}

	// RULE 2: Blockage Detection (TREND)
	if len(state.LastDistances) >= 3 {
		l1 := state.LastDistances[len(state.LastDistances)-3]
		l2 := state.LastDistances[len(state.LastDistances)-2]
		l3 := state.LastDistances[len(state.LastDistances)-1]

		// Jarak makin ngecil (air makin naik) berturut-turut + tidak hujan
		if l1 > l2 && l2 > l3 && !rainDetected {
			return utils.Blockage, true
		}
	}

	// RULE 3: Informasi Hujan Biasa (Bukan Alert Kritikal)
	return utils.RainInfo, rainDetected
}
