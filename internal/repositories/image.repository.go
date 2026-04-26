package repositories

import (
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v3"
)

func (r *Repository) UploadImage(c fiber.Ctx, src multipart.File) (string, *fiber.Error) {
	result, err := r.Cloudinary.Cloudinary.Upload.Upload(c, src, uploader.UploadParams{Folder: "COMPRO NEED"})
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return result.SecureURL, nil
}
