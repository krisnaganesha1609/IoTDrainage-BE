package repositories

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v3"
)

// UpsertDeviceInfo stores / updates device maintenance info (IP, location)
// in Firestore. Also performs auto-registration if the document does not exist.
func (r *Repository) UpsertDeviceInfo(deviceID, ipAddress, location string, deviceTimestamp int64) *fiber.Error {
	ctx := context.Background()

	t := time.Unix(deviceTimestamp, 0).UTC()

	_, err := r.Firebase.Firestore.Collection("devices").Doc(deviceID).Set(ctx, map[string]interface{}{
		"device_id":      deviceID,
		"ip_address":     ipAddress,
		"location":       location,
		"last_seen_boot": t,
		// Set status ONLINE on boot; watchdog will flip to OFFLINE later.
		"status": "ONLINE",
	}, firestore.MergeAll)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan device-info ke Firestore: "+err.Error())
	}

	return nil
}

// UpdateDeviceHeartbeat updates the last-seen telemetry timestamp and
// next_wakeup_sec so the watchdog cron can detect offline devices.
func (r *Repository) UpdateDeviceHeartbeat(deviceID string, lastTimestamp, nextWakeupSec int64) *fiber.Error {
	ctx := context.Background()

	_, err := r.Firebase.Firestore.Collection("devices").Doc(deviceID).Set(ctx, map[string]interface{}{
		"last_timestamp":  lastTimestamp,
		"next_wakeup_sec": nextWakeupSec,
		"status":          "ONLINE",
	}, firestore.MergeAll)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal update heartbeat ke Firestore: "+err.Error())
	}

	return nil
}
