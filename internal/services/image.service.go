package services

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/entities"
)

func (s *Service) UploadImage(c fiber.Ctx, file *multipart.FileHeader) (string, *fiber.Error) {
	src, err := file.Open()
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "cannot open file")
	}
	defer src.Close()
	secureURL, errImg := s.Repo.UploadImage(c, src)
	if errImg != nil {
		return "", errImg
	}
	return s.Repo.InsertImageToInflux(secureURL)
}

func (s *Service) GetLatestImage() (*entities.ImageEntity, *fiber.Error) {
	result, err := s.Repo.GetLatestImage()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var imageURL string
	var lastUpdated int64
	hasData := false

	for result.Next() {

		record := result.Record()
		if record == nil {
			continue
		}

		if val, ok := record.ValueByKey("image_url").(string); ok {
			imageURL = val
			lastUpdated = record.Time().Unix()
			hasData = true
		}
	}

	if result.Err() != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Err().Error())
	}

	if !hasData {
		return nil, fiber.NewError(fiber.StatusNotFound, "No image found")
	}

	return &entities.ImageEntity{
		URL:         imageURL,
		LastUpdated: lastUpdated,
	}, nil
}
