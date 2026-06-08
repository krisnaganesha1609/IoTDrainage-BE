package services

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/iterator"
)

const watchdogToleranceSec = 5 * 60 // 5-minute tolerance per spec

// RunWatchdog starts a goroutine that periodically checks every known device
// in Firestore. If a device has not reported within its expected wakeup
// window + 5-minute tolerance, its Firestore status is flipped to "OFFLINE".
//
// Formula (per spec):
//
//	now > last_timestamp + next_wakeup_sec + 300  =>  OFFLINE
func (s *Service) RunWatchdog() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("[Watchdog] Started — checking every 1 minute")

	for range ticker.C {
		s.checkOfflineDevices()
	}
}

func (s *Service) checkOfflineDevices() {
	ctx := context.Background()
	now := time.Now().Unix()

	iter := s.Firebase.Firestore.Collection("devices").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("[Watchdog] Error iterating devices: %v", err)
			break
		}

		data := doc.Data()
		deviceID := doc.Ref.ID

		// Skip devices that have already been marked offline to avoid spam writes.
		if status, ok := data["status"].(string); ok && status == "OFFLINE" {
			continue
		}

		lastTS, okTS := data["last_timestamp"].(int64)
		nextWakeup, okWake := data["next_wakeup_sec"].(int64)

		if !okTS || !okWake || nextWakeup == 0 {
			// Not enough data to evaluate; device may not have sent telemetry yet.
			continue
		}

		deadline := lastTS + nextWakeup + watchdogToleranceSec
		if now > deadline {
			log.Printf("[Watchdog] Device %s dinyatakan OFFLINE (deadline %d < now %d)", deviceID, deadline, now)
			if err := s.Repo.SetDeviceOffline(deviceID); err != nil {
				log.Printf("[Watchdog] Gagal set OFFLINE device %s: %v", deviceID, err)
			}
		}
	}
}
