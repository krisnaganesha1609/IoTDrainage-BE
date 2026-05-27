package services

import (
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/entities"
)

func (s *Service) UploadImage(c fiber.Ctx, file *multipart.FileHeader, deviceID string) (string, *fiber.Error) {
	src, err := file.Open()
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "cannot open file")
	}
	defer src.Close()

	secureURL, errImg := s.Repo.UploadImage(c, src)
	if errImg != nil {
		return "", errImg
	}

	return s.Repo.InsertImageToFirestore(deviceID, secureURL)
}

func (s *Service) GetLatestImage() (*entities.ImageEntity, *fiber.Error) {
	data, errRef := s.Repo.GetLatestImageFromFirestore()
	if errRef != nil {
		return nil, errRef
	}

	imageURL, okURL := data["image_url"].(string)
	lastUpdated, okTime := data["last_updated"].(int64)

	if !okURL {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Format image_url rusak di database")
	}
	if !okTime {
		if t, ok := data["created_at"].(time.Time); ok {
			lastUpdated = t.Unix()
		}
	}

	return &entities.ImageEntity{
		URL:         imageURL,
		LastUpdated: lastUpdated,
	}, nil
}
