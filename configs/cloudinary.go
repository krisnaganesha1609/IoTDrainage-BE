package configs

import (
	"github.com/gofiber/fiber/v3/log"

	"github.com/cloudinary/cloudinary-go/v2"
)

type CloudinaryConfig struct {
	Cloudinary *cloudinary.Cloudinary
}

func InitCloudinary(url string) (*CloudinaryConfig, error) {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		log.Fatalf("Failed to initialize cloudinary: %v", err)
	}
	return &CloudinaryConfig{Cloudinary: cld}, nil
}
