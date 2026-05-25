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

// Struct penampung data yang di-unmarshal dari Firestore dokumen
type FirestoreDeviceDoc struct {
	MobileDevices map[string]struct {
		FCMToken  string    `firestore:"fcm_token"`
		UpdatedAt time.Time `firestore:"updated_at"`
	} `firestore:"mobile_devices"`
}

func (s *Service) ProcessSensorData(request requests.SensorDataRequest) *fiber.Error {
	// 1. Selalu simpan data mentah dari sensor ke InfluxDB untuk keperluan historis/analis
	if err := s.Repo.InsertSensorData(request.DeviceID, request.Location, request.WaterDistance, request.RainDetected, request.RainIntensity); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// 2. Jalankan Complex Event Processing (CEP) Engine di RAM untuk cek Rule & Tren
	alertType, triggered := core.ProcessEvent(request.DeviceID, request.WaterDistance, request.RainDetected)

	// CASE A: Tidak Ada Alert / Kondisi Aman
	if !triggered {
		forwarding := entities.WebsocketMessage{
			Type: "sensor",
			Data: map[string]interface{}{
				"payload":   request,
				"timestamp": time.Now().Unix(),
			},
		}

		// Throttle Broadcast Websocket data normal (misal tiap 20 detik sekali) agar tidak spamming browser
		if time.Since(utils.DeviceStates[request.DeviceID].LastTimes[4]) > 20*time.Second {
			s.Broadcast(forwarding)
		}
	}

	// CASE B: Terjadi Trigger Kondisi Bahaya (HIGH_WATER atau BLOCKAGE)
	if triggered {
		// 1. Catat kejadian alert ke database repository (InfluxDB/Log)
		if err := s.Repo.InsertAlert(request.DeviceID, string(alertType)); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// 2. Kirim alert instan via Websocket ke Dashboard Web UI yang sedang stand-by
		msg := entities.WebsocketMessage{
			Type: "alert",
			Data: map[string]interface{}{
				"alert":     string(alertType),
				"payload":   request,
				"timestamp": time.Now().Unix(),
			},
		}
		s.Broadcast(msg)

		// 3. PUSH NOTIFICATION VIA FIRESTORE & FCM MULTICAST (ASYNC)
		go func(ctx context.Context, deviceID string, aType utils.AlertType) {
			// Query dokumen alat IoT dari koleksi "devices" di Firestore
			dsnap, err := s.Firebase.Firestore.Collection("devices").Doc(deviceID).Get(ctx)
			if err != nil {
				log.Printf("[FCM] Gagal mengambil data dokumen dari Firestore: %v", err)
				return
			}

			if !dsnap.Exists() {
				log.Printf("[FCM] Dokumen device %s tidak ditemukan di Firestore, skip notifikasi.", deviceID)
				return
			}

			// Unmarshal data sub-koleksi map mobile_devices
			var deviceDoc FirestoreDeviceDoc
			if err := dsnap.DataTo(&deviceDoc); err != nil {
				log.Printf("[FCM] Gagal unmarshal data Firestore: %v", err)
				return
			}

			// Ekstrak seluruh token HP yang terdaftar untuk dikirim bersamaan
			var tokens []string
			for _, mobile := range deviceDoc.MobileDevices {
				if mobile.FCMToken != "" {
					tokens = append(tokens, mobile.FCMToken)
				}
			}

			// Jika ada HP yang terhubung, tembak menggunakan Multicast (1x hit API Firebase)
			if len(tokens) > 0 {
				multicastMessage := &messaging.MulticastMessage{
					Tokens: tokens,
					Notification: &messaging.Notification{
						Title: "⚠️ Peringatan Sistem Drainase!",
						Body:  fmt.Sprintf("Alat %s mendeteksi indikasi bahaya: %s", deviceID, string(aType)),
					},
					Data: map[string]string{
						"device_id":  deviceID,
						"alert_type": string(aType),
					},
				}

				if s.Firebase.FCM != nil {
					br, sendErr := s.Firebase.FCM.SendEachForMulticast(ctx, multicastMessage)
					if sendErr != nil {
						log.Printf("[FCM] Gagal mengirim Multicast Notifikasi: %v", sendErr)
					} else {
						log.Printf("[FCM] Notifikasi Sukses dikirim ke %d HP (Gagal: %d)", br.SuccessCount, br.FailureCount)
					}
				}
			}
		}(context.Background(), request.DeviceID, alertType)

		// 4. Update LastAlertTime di RAM utilitas lokal agar cooldown 2 menit berjalan
		utils.DeviceStates[request.DeviceID].LastAlertTime = time.Now()

		// 5. Cadangkan juga data LastAlertTime ke Firestore secara async agar persisten
		go func(ctx context.Context, deviceID string) {
			_, _ = s.Firebase.Firestore.Collection("devices").Doc(deviceID).Set(ctx, map[string]interface{}{
				"last_alert_time": time.Now(),
			}, firestore.MergeAll)
		}(context.Background(), request.DeviceID)
	}

	return nil
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
		sensorHistory := entities.SensorHistory{
			Time:          record.Time(),
			WaterDistance: record.ValueByKey("water_distance").(float64),
			RainDetected:  record.ValueByKey("rain_detected").(bool),
			RainIntensity: record.ValueByKey("rain_intensity").(float64),
		}
		sensorData = append(sensorData, sensorHistory)
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
		sensorHistory := entities.SensorHistory{
			Time:          record.Time(),
			WaterDistance: record.ValueByKey("water_distance").(float64),
			RainDetected:  record.ValueByKey("rain_detected").(bool),
			RainIntensity: record.ValueByKey("rain_intensity").(float64),
		}
		sensorData = append(sensorData, sensorHistory)
	}
	return sensorData, nil
}
