package repositories

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InsertSensorData writes a telemetry point to InfluxDB.
// The timestamp comes from the device (UTC Unix), not server time,
// so that InfluxDB accurately reflects when the reading was taken.
func (r *Repository) InsertSensorData(
	deviceID string,
	waterDistance float64,
	waterLevelCm float64,
	status string,
	rainDetected bool,
	sensorFlag string,
	nextWakeupSec int64,
	deviceTimestamp int64,
) error {
	t := time.Unix(deviceTimestamp, 0).UTC()

	point := influxdb2.NewPoint(
		"drainage_monitoring",
		map[string]string{
			"device_id": deviceID,
			"status":    status,
		},
		map[string]interface{}{
			"water_distance":  waterDistance,
			"water_level_cm":  waterLevelCm,
			"rain_detected":   rainDetected,
			"sensor_flag":     sensorFlag,
			"next_wakeup_sec": nextWakeupSec,
		},
		t,
	)

	return r.Influx.WriteAPI.WritePoint(context.Background(), point)
}

func (r *Repository) GetSensorHistory(deviceID, startDate, endDate string) (*api.QueryTableResult, error) {
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: %s, stop: %s)
  |> filter(fn: (r) => r._measurement == "drainage_monitoring")
  |> filter(fn: (r) => r.device_id == "%s")
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> sort(columns: ["_time"], desc: true)
`, r.Bucket, startDate, endDate, deviceID)

	return r.Influx.QueryAPI.Query(context.Background(), query)
}

func (r *Repository) GetLatestSensorData(deviceID string) (*api.QueryTableResult, error) {
	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "drainage_monitoring")
  |> filter(fn: (r) => r.device_id == "%s")
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: 1)
`, r.Bucket, deviceID)

	return r.Influx.QueryAPI.Query(context.Background(), query)
}
