package repositories

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InsertSensorLog writes a system-log entry from the IoT device into InfluxDB.
func (r *Repository) InsertSensorLog(
	deviceID string,
	level string,
	message string,
	wakeReason string,
	resetReason string,
	activeTimeMs int64,
	wifiRssiDbm int64,
	networkFailures int64,
	freeHeapBytes int64,
	deviceTimestamp int64,
) error {
	t := time.Unix(deviceTimestamp, 0).UTC()

	point := influxdb2.NewPoint(
		"device_logs",
		map[string]string{
			"device_id": deviceID,
			"level":     level,
		},
		map[string]interface{}{
			"message":          message,
			"wake_reason":      wakeReason,
			"reset_reason":     resetReason,
			"active_time_ms":   activeTimeMs,
			"wifi_rssi_dbm":    wifiRssiDbm,
			"network_failures": networkFailures,
			"free_heap_bytes":  freeHeapBytes,
		},
		t,
	)

	return r.Influx.WriteAPI.WritePoint(context.Background(), point)
}

// GetRecentLogs fetches device logs from the last 7 days from InfluxDB.
func (r *Repository) GetRecentLogs(deviceID string) (*api.QueryTableResult, error) {
	deviceFilter := ""
	if deviceID != "" {
		deviceFilter = fmt.Sprintf(`|> filter(fn: (r) => r.device_id == "%s")`, deviceID)
	}

	query := fmt.Sprintf(`
from(bucket: "%s")
  |> range(start: -7d)
  |> filter(fn: (r) => r._measurement == "device_logs")
  %s
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> sort(columns: ["_time"], desc: true)
`, r.Bucket, deviceFilter)

	return r.Influx.QueryAPI.Query(context.Background(), query)
}
