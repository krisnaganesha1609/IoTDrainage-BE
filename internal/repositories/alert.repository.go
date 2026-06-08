package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (r *Repository) InsertAlert(deviceID, alertType string) error {
	point := influxdb2.NewPoint(
		"alerts",
		map[string]string{
			"device_id": deviceID,
			"type":      alertType,
		},
		map[string]interface{}{
			"value": 1,
		},
		time.Now(),
	)
	return r.Influx.WriteAPI.WritePoint(context.Background(), point)
}

// SetDeviceOffline marks a device as OFFLINE in Firestore.
// Called by the watchdog cron when the device has not reported within
// its expected wakeup window + 5-minute tolerance.
func (r *Repository) SetDeviceOffline(deviceID string) error {
	ctx := context.Background()

	_, err := r.Firebase.Firestore.Collection("devices").Doc(deviceID).Set(ctx, map[string]interface{}{
		"status":        "OFFLINE",
		"offline_since": time.Now().UTC(),
	}, firestore.MergeAll)

	if err != nil {
		return fmt.Errorf("gagal set device %s menjadi OFFLINE: %w", deviceID, err)
	}
	return nil
}
