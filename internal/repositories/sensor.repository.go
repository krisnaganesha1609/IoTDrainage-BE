package repositories

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func (r *Repository) InsertSensorData(deviceID, location string, waterDistance float64, rainDetected bool, rainIntensity float64) error {
	point := influxdb2.NewPoint(
		"drainage_monitoring",
		map[string]string{
			"device_id": deviceID,
			"location":  location,
		},
		map[string]interface{}{
			"water_distance": waterDistance,
			"rain_detected":  rainDetected,
			"rain_intensity": rainIntensity,
		},
		time.Now(),
	)

	return r.Influx.WriteAPI.WritePoint(context.Background(), point)
}

func (r *Repository) GetSensorHistory(deviceID, startDate, endDate string) (*api.QueryTableResult, error) {
	query := `
	from(bucket: "iot_drainage")
	|> range(start: ` + startDate + `, stop: ` + endDate + `)
	|> filter(fn: (r) => r._measurement == "drainage_monitoring")
	|> filter(fn: (r) => r.device_id == "` + deviceID + `")
	|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	|> sort(columns: ["_time"], desc: true)
	`
	result, err := r.Influx.QueryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetLatestSensorData(deviceID string) (*api.QueryTableResult, error) {
	query := `
	from(bucket: "iot_drainage")
	|> range(start: -1h)
	|> filter(fn: (r) => r._measurement == "drainage_monitoring")
	|> filter(fn: (r) => r.device_id == "` + deviceID + `")
	|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	|> sort(columns: ["_time"], desc: true)
	`
	result, err := r.Influx.QueryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
