package services

import (
	"farmish/internal/models"
	"farmish/internal/repository"
	"time"

	"github.com/google/uuid"
)

type FarmService struct {
	repo *repository.FarmRepository
}

func NewFarmService(repo *repository.FarmRepository) *FarmService {
	return &FarmService{repo: repo}
}

func (s *FarmService) CreateFarm(farm *models.Farm) error {
	farm.ID = uuid.New()
	farm.CreatedAt = time.Now()
	return s.repo.CreateFarm(farm)
}

func (s *FarmService) GetFarmByID(farmID uuid.UUID) (*models.Farm, error) {
	return s.repo.GetFarmByID(farmID)
}

func (s *FarmService) GetAllFarms() ([]models.Farm, error) {
	return s.repo.GetAllFarms()
}

func (s *FarmService) UpdateFarm(farm *models.UpdateFarmRequest) error {
	return s.repo.UpdateFarm(farm)
}

func (s *FarmService) DeleteFarm(farmID uuid.UUID) error {
	return s.repo.DeleteFarm(farmID)
}
