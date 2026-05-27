package repositories

import (
	"github.com/krisnaganesha1609/IoTDrainage-BE/configs"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

type Repository struct {
	Influx     *configs.InfluxDB
	Cloudinary *configs.CloudinaryConfig
	Firebase   *utils.FirebaseServices
}

func InitializeRepository(influx *configs.InfluxDB, cloudinary *configs.CloudinaryConfig, firebase *utils.FirebaseServices) *Repository {
	return &Repository{
		Influx:     influx,
		Cloudinary: cloudinary,
		Firebase:   firebase,
	}
}
