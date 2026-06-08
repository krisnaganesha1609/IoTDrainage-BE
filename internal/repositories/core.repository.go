package repositories

import (
	"github.com/krisnaganesha1609/IoTDrainage-BE/configs"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

type Repository struct {
	Influx     *configs.InfluxDB
	Cloudinary *configs.CloudinaryConfig
	Firebase   *utils.FirebaseServices
	// Bucket name is stored here so Flux queries don't hardcode it.
	Bucket string
}

func InitializeRepository(influx *configs.InfluxDB, cloudinary *configs.CloudinaryConfig, firebase *utils.FirebaseServices, bucket string) *Repository {
	return &Repository{
		Influx:     influx,
		Cloudinary: cloudinary,
		Firebase:   firebase,
		Bucket:     bucket,
	}
}
