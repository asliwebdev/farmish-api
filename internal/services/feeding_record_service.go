package services

import (
	"errors"
	"farmish/internal/models"
	"farmish/internal/repository"

	"github.com/google/uuid"
)

type FeedingRecordService struct {
	feedingRecordRepo *repository.FeedingRecordRepository
	animalRepo        *repository.AnimalRepository
	foodRepo          *repository.FoodRepository
}

func NewFeedingRecordService(
	feedingRecordRepo *repository.FeedingRecordRepository,
	animalRepo *repository.AnimalRepository,
	foodRepo *repository.FoodRepository,
) *FeedingRecordService {
	return &FeedingRecordService{
		feedingRecordRepo: feedingRecordRepo,
		animalRepo:        animalRepo,
		foodRepo:          foodRepo,
	}
}

var (
	ErrAnimalNotFound       = errors.New("animal not found")
	ErrFoodNotFound         = errors.New("food not found")
	ErrInsufficientQuantity = errors.New("insufficient food quantity")
)

func (s *FeedingRecordService) CreateFeedingRecord(record *models.FeedingRecordWithoutTime) error {
	animal, err := s.animalRepo.GetAnimalByID(record.AnimalID)
	if err != nil {
		return err
	} else if animal == nil {
		return ErrAnimalNotFound
	}

	food, err := s.foodRepo.GetFoodByID(record.FoodID)
	if err != nil {
		return err
	} else if food == nil {
		return ErrFoodNotFound
	}

	if food.Quantity < record.Quantity {
		return ErrInsufficientQuantity
	}

	newQuantity := food.Quantity - record.Quantity

	record.ID = uuid.New()

	return s.feedingRecordRepo.CreateFeedingRecord(record, newQuantity)
}

func (s *FeedingRecordService) GetFeedingRecordByID(id uuid.UUID) (*models.FeedingRecordDetailed, error) {
	return s.feedingRecordRepo.GetFeedingRecordByID(id)
}

func (s *FeedingRecordService) GetFeedingRecordsByAnimalID(animalID uuid.UUID) ([]models.FeedingRecordDetailed, error) {
	return s.feedingRecordRepo.GetFeedingRecordsByAnimalID(animalID)
}

func (s *FeedingRecordService) UpdateFeedingRecord(record *models.FeedingRecordWithoutTime) error {
	return s.feedingRecordRepo.UpdateFeedingRecord(record)
}

func (s *FeedingRecordService) DeleteFeedingRecord(id uuid.UUID) error {
	return s.feedingRecordRepo.DeleteFeedingRecord(id)
}
