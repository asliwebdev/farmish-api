package services

import (
	"farmish/internal/models"
	"farmish/internal/repository"

	"github.com/google/uuid"
)

type FoodService struct {
	FoodRepo *repository.FoodRepository
}

func NewFoodService(repo *repository.FoodRepository) *FoodService {
	return &FoodService{FoodRepo: repo}
}

func (s *FoodService) AddFoodToWarehouse(food *models.FoodWithoutTime) error {
	food.ID = uuid.New()
	return s.FoodRepo.CreateFood(food)
}

func (s *FoodService) GetFoodsByFarm(farmID uuid.UUID) ([]models.Food, error) {
	return s.FoodRepo.GetAllFoods(farmID)
}

func (s *FoodService) GetFoodByID(foodID uuid.UUID) (*models.Food, error) {
	return s.FoodRepo.GetFoodByID(foodID)
}

func (s *FoodService) UpdateFood(food *models.UpdateFoodReq) error {
	return s.FoodRepo.UpdateFood(food)
}

func (s *FoodService) RemoveWarehouseFood(foodID uuid.UUID) error {
	return s.FoodRepo.DeleteFood(foodID)
}
