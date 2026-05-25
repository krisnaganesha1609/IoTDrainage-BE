package services

import (
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/repositories"
	"github.com/krisnaganesha1609/IoTDrainage-BE/utils"
)

type Service struct {
	Repo     *repositories.Repository
	Firebase *utils.FirebaseServices
}

func InitializeService(repo *repositories.Repository, firebase *utils.FirebaseServices) *Service {
	return &Service{
		Repo:     repo,
		Firebase: firebase,
	}
}
