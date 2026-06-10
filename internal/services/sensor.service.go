package services

import (
	"context"
	"fmt"
	"log"
	"time"

	firestore "cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/core"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/entities"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

// FirestoreDeviceDoc maps the Firestore device document for FCM token retrieval.
type FirestoreDeviceDoc struct {
	MobileDevices map[string]struct {
		FCMToken  string    `firestore:"fcm_token"`
		UpdatedAt time.Time `firestore:"updated_at"`
	} `firestore:"mobile_devices"`
}

// ProcessSensorData is the central handler for incoming MQTT telemetry.
//
// Architecture (Smart Publisher):
//   - IoT's `status` field is Source of Truth (NORMAL/WASPADA/BAHAYA).
//   - FCM is triggered when status == "BAHAYA" (from IoT) OR BLOCKAGE (from CEP).
//   - Backend does NOT recalculate the flood status; it only stores and forwards.
func (s *Service) ProcessSensorData(request requests.SensorDataRequest) *fiber.Error {
	if request.DeviceID == "" {
		log.Println("[WARN] Payload sensor ditolak: device_id kosong")
		return fiber.NewError(fiber.StatusBadRequest, "device_id is required")
	}

	// 1. Persist raw telemetry to InfluxDB using device-provided UTC timestamp.
	if err := s.Repo.InsertSensorData(
		request.DeviceID,
		request.WaterDistance,
		request.WaterLevelCm,
		request.Status,
		request.RainDetected,
		request.SensorFlag,
		request.NextWakeupSec,
		request.Timestamp,
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// 2. Update device heartbeat in Firestore for the watchdog cron.
	go func() {
		if err := s.Repo.UpdateDeviceHeartbeat(request.DeviceID, request.Timestamp, request.NextWakeupSec); err != nil {
			log.Printf("[Heartbeat] Gagal update heartbeat device %s: %v", request.DeviceID, err)
		}
	}()

	// 3. Broadcast live telemetry via WebSocket (throttled to 20 s for non-alert data).
	wsMsg := entities.WebsocketMessage{
		Type: "sensor",
		Data: map[string]interface{}{
			"payload":   request,
			"timestamp": time.Now().Unix(),
		},
	}

	state := utils.GetOrCreateDeviceState(request.DeviceID)
	utils.DeviceStatesMu.RLock()
	lastBroadcast := state.LastTimes[4]
	utils.DeviceStatesMu.RUnlock()

	if request.Status == "BAHAYA" || time.Since(lastBroadcast) > 20*time.Second {
		s.Broadcast(wsMsg)
		utils.DeviceStatesMu.Lock()
		state.LastTimes[4] = time.Now()
		utils.DeviceStatesMu.Unlock()
	}

	// 4a. FCM trigger: IoT reports BAHAYA status (Source of Truth).
	if request.Status == "BAHAYA" {
		// utils.DeviceStatesMu.RLock()
		// cooldownOK := time.Since(state.LastAlertTime) >= 2*time.Minute
		// utils.DeviceStatesMu.RUnlock()

		// if cooldownOK {
		if err := s.Repo.InsertAlert(request.DeviceID, "BAHAYA"); err != nil {
			log.Printf("[Alert] Gagal insert alert BAHAYA: %v", err)
		}

		alertWsMsg := entities.WebsocketMessage{
			Type: "alert",
			Data: map[string]interface{}{
				"alert":     "BAHAYA",
				"payload":   request,
				"timestamp": time.Now().Unix(),
			},
		}
		s.Broadcast(alertWsMsg)

		go s.sendFCMAlert(request.DeviceID, "BAHAYA", request)

		utils.DeviceStatesMu.Lock()
		state.LastAlertTime = time.Now()
		utils.DeviceStatesMu.Unlock()

		// Persist last alert time to Firestore for cross-restart durability.
		go func(deviceID string) {
			ctx := context.Background()
			_, _ = s.Firebase.Firestore.Collection("devices").Doc(deviceID).Set(ctx,
				map[string]interface{}{"last_alert_time": time.Now()},
				firestore.MergeAll,
			)
		}(request.DeviceID)
		// }
	}

	// 4b. FCM trigger: CEP detects BLOCKAGE (BE-side, not covered by IoT status).
	if alertType, triggered := core.ProcessEvent(request.DeviceID, request.WaterDistance, request.RainDetected); triggered {
		if err := s.Repo.InsertAlert(request.DeviceID, string(alertType)); err != nil {
			log.Printf("[Alert] Gagal insert alert %s: %v", alertType, err)
		}

		blockageWsMsg := entities.WebsocketMessage{
			Type: "alert",
			Data: map[string]interface{}{
				"alert":     string(alertType),
				"payload":   request,
				"timestamp": time.Now().Unix(),
			},
		}
		s.Broadcast(blockageWsMsg)

		go s.sendFCMAlert(request.DeviceID, string(alertType), request)

		utils.DeviceStatesMu.Lock()
		state.LastAlertTime = time.Now()
		utils.DeviceStatesMu.Unlock()
	}

	return nil
}

// sendFCMAlert fetches FCM tokens from Firestore and sends a multicast push notification.
func (s *Service) sendFCMAlert(deviceID, alertType string, request requests.SensorDataRequest) {
	ctx := context.Background()

	dsnap, err := s.Firebase.Firestore.Collection("devices").Doc(deviceID).Get(ctx)
	if err != nil || !dsnap.Exists() {
		log.Printf("[FCM] Dokumen device %s tidak ditemukan, skip notifikasi.", deviceID)
		return
	}

	var deviceDoc FirestoreDeviceDoc
	if err := dsnap.DataTo(&deviceDoc); err != nil {
		log.Printf("[FCM] Gagal unmarshal data Firestore: %v", err)
		return
	}

	var tokens []string
	for _, mobile := range deviceDoc.MobileDevices {
		if mobile.FCMToken != "" {
			tokens = append(tokens, mobile.FCMToken)
		}
	}

	if len(tokens) == 0 || s.Firebase.FCM == nil {
		return
	}

	var title string
	var body string

	switch request.Status {
	case "NORMAL":
		title = "✅ Status Drainase: Normal"
		body = fmt.Sprintf(
			"%s: tinggi air %.1f cm. Sistem beroperasi normal.",
			deviceID,
			request.WaterLevelCm,
		)

	case "WASPADA":
		title = "⚠️ Status Drainase: Waspada"
		body = fmt.Sprintf(
			"%s: tinggi air %.1f cm. Terjadi kenaikan muka air, pemantauan disarankan.",
			deviceID,
			request.WaterLevelCm,
		)

	case "BAHAYA":
		title = "🚨 Status Drainase: Bahaya"
		body = fmt.Sprintf(
			"%s: tinggi air %.1f cm. Tinggi air mencapai batas kritis, waspada potensi banjir!",
			deviceID,
			request.WaterLevelCm,
		)

	default:
		title = fmt.Sprintf("📡 Update Sensor Drainase: %s", alertType)
		body = fmt.Sprintf(
			"%s: tinggi air %.1f cm. Status: %s.",
			deviceID,
			request.WaterLevelCm,
			request.Status,
		)
	}

	multicastMessage := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"device_id":  deviceID,
			"alert_type": alertType,
		},
		Android: &messaging.AndroidConfig{Priority: "high"},
	}

	br, sendErr := s.Firebase.FCM.SendEachForMulticast(ctx, multicastMessage)
	if sendErr != nil {
		log.Printf("[FCM] Gagal mengirim multicast: %v", sendErr)
	} else {
		log.Printf("[FCM] Notifikasi dikirim ke %d HP (Gagal: %d)", br.SuccessCount, br.FailureCount)
	}
}

func (s *Service) GetSensorHistory(deviceID, startDate, endDate string) ([]entities.SensorHistory, *fiber.Error) {
	result, err := s.Repo.GetSensorHistory(deviceID, startDate, endDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer result.Close()

	var sensorData []entities.SensorHistory
	for result.Next() {
		record := result.Record()
		sensorData = append(sensorData, entities.SensorHistory{
			Time:          record.Time(),
			DeviceID:      utils.SafeString(record.ValueByKey("device_id")),
			WaterDistance: utils.SafeFloat64(record.ValueByKey("water_distance")),
			WaterLevelCm:  utils.SafeFloat64(record.ValueByKey("water_level_cm")),
			Status:        utils.SafeString(record.ValueByKey("status")),
			RainDetected:  utils.SafeBool(record.ValueByKey("rain_detected")),
			SensorFlag:    utils.SafeString(record.ValueByKey("sensor_flag")),
			NextWakeupSec: utils.SafeInt64(record.ValueByKey("next_wakeup_sec")),
		})
	}
	return sensorData, nil
}

func (s *Service) GetLatestSensorData(deviceID string) ([]entities.SensorHistory, *fiber.Error) {
	result, err := s.Repo.GetLatestSensorData(deviceID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer result.Close()

	var sensorData []entities.SensorHistory
	for result.Next() {
		record := result.Record()
		sensorData = append(sensorData, entities.SensorHistory{
			Time:          record.Time(),
			DeviceID:      utils.SafeString(record.ValueByKey("device_id")),
			WaterDistance: utils.SafeFloat64(record.ValueByKey("water_distance")),
			WaterLevelCm:  utils.SafeFloat64(record.ValueByKey("water_level_cm")),
			Status:        utils.SafeString(record.ValueByKey("status")),
			RainDetected:  utils.SafeBool(record.ValueByKey("rain_detected")),
			SensorFlag:    utils.SafeString(record.ValueByKey("sensor_flag")),
			NextWakeupSec: utils.SafeInt64(record.ValueByKey("next_wakeup_sec")),
		})
	}
	return sensorData, nil
}
