package services

import (
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/repositories"
)

type Service struct {
	Repo *repositories.Repository
}

func InitializeService(repo *repositories.Repository) *Service {
	return &Service{Repo: repo}
}
