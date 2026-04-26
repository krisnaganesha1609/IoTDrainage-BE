package repositories

import "github.com/krisnaganesha1609/IoTDrainage-BE/configs"

type Repository struct {
	Influx     *configs.InfluxDB
	Cloudinary *configs.CloudinaryConfig
}

func InitializeRepository(influx *configs.InfluxDB, cloudinary *configs.CloudinaryConfig) *Repository {
	return &Repository{
		Influx:     influx,
		Cloudinary: cloudinary,
	}
}
