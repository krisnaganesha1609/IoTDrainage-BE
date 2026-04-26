package services

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) UploadImage(c fiber.Ctx, file *multipart.FileHeader) (string, *fiber.Error) {
	src, err := file.Open()
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "cannot open file")
	}
	defer src.Close()
	return s.Repo.UploadImage(c, src)
}
