package repositories

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v3"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func (r *Repository) UploadImage(c fiber.Ctx, src multipart.File) (string, *fiber.Error) {
	result, err := r.Cloudinary.Cloudinary.Upload.Upload(c, src, uploader.UploadParams{Folder: "COMPRO NEED"})
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return result.SecureURL, nil
}

func (r *Repository) InsertImageToInflux(imageURL string) (string, *fiber.Error) {
	point := influxdb2.NewPoint(
		"drainage_images",
		map[string]string{},
		map[string]interface{}{
			"image_url": imageURL,
		},
		time.Now(),
	)
	if err := r.Influx.WriteAPI.WritePoint(context.Background(), point); err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return imageURL, nil
}

func (r *Repository) GetLatestImage() (*api.QueryTableResult, *fiber.Error) {
	query := `
    from(bucket: "iot_drainage")
    |> range(start: 0)
    |> filter(fn: (r) => r._measurement == "drainage_images")
    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
    |> sort(columns: ["_time"], desc: true)
    |> limit(n:1)
    `
	result, err := r.Influx.QueryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return result, nil
}
