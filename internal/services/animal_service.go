package services

import (
	"errors"
	"farmish/internal/models"
	"farmish/internal/repository"

	"github.com/google/uuid"
)

type AnimalService struct {
	Repo *repository.AnimalRepository
}

func NewAnimalService(repo *repository.AnimalRepository) *AnimalService {
	return &AnimalService{Repo: repo}
}

var ErrNegativeWeight = errors.New("weight must be greater than 0")

func (s *AnimalService) CreateAnimal(animal *models.AnimalWithoutTime) error {
	animal.ID = uuid.New()

	if animal.Weight <= 0 {
		return ErrNegativeWeight
	}

	return s.Repo.CreateAnimal(animal)
}

func (s *AnimalService) GetAnimalByID(animalID uuid.UUID) (*models.Animal, error) {
	return s.Repo.GetAnimalByID(animalID)
}

func (s *AnimalService) GetAnimalsByFarmID(farmID uuid.UUID) ([]*models.Animal, error) {
	return s.Repo.GetAnimalsByFarmID(farmID)
}

func (s *AnimalService) UpdateAnimal(animal *models.UpdateAnimalReq) error {
	if animal.Weight <= 0 {
		return ErrNegativeWeight
	}

	return s.Repo.UpdateAnimal(animal)
}

func (s *AnimalService) DeleteAnimal(animalID uuid.UUID) error {
	return s.Repo.DeleteAnimal(animalID)
}
