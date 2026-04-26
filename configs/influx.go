package configs

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDB struct {
	InfluxClient influxdb2.Client
	WriteAPI     api.WriteAPIBlocking
	QueryAPI     api.QueryAPI
}

func InitInfluxDB(influxURL string, influxToken string, influxOrg string, influxBucket string) (*InfluxDB, error) {
	influxClient := influxdb2.NewClient(influxURL, influxToken)
	writeAPI := influxClient.WriteAPIBlocking(influxOrg, influxBucket)
	queryAPI := influxClient.QueryAPI(influxOrg)

	return &InfluxDB{
		InfluxClient: influxClient,
		WriteAPI:     writeAPI,
		QueryAPI:     queryAPI,
	}, nil
}
