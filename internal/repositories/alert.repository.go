package repositories

import (
	"context"
	"time"

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
