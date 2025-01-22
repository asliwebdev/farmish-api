package services

import (
	"errors"
	"farmish/internal/models"
	"farmish/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

type MedicineService struct {
	repo repository.MedicineRepository
}

func NewMedicineService(repo repository.MedicineRepository) *MedicineService {
	return &MedicineService{repo: repo}
}

var (
	ErrQuantityLessThanThreshold = errors.New("quantity cannot be less than the minimum threshold")
	ErrNotExist                  = errors.New("medicine with this ID not found")
)

func (s *MedicineService) CreateMedicine(medicine *models.MedicineWithoutTime) error {
	medicine.ID = uuid.New()
	if medicine.Quantity < medicine.MinThreshold {
		return ErrQuantityLessThanThreshold
	}

	return s.repo.CreateMedicine(medicine)
}

func (s *MedicineService) GetAllMedicines(farmID uuid.UUID) ([]models.Medicine, error) {
	return s.repo.GetAllMedicines(farmID)
}

func (s *MedicineService) GetMedicineByID(id uuid.UUID) (*models.Medicine, error) {
	return s.repo.GetMedicineByID(id)
}

func (s *MedicineService) UpdateMedicine(medicine *models.MedicineWithoutTime) error {
	existing, err := s.repo.GetMedicineByID(medicine.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing medicine: %v", err)
	}
	if existing == nil {
		return ErrNotExist
	}

	if medicine.Quantity < medicine.MinThreshold {
		return ErrQuantityLessThanThreshold
	}

	return s.repo.UpdateMedicine(medicine)
}

func (s *MedicineService) DeleteMedicine(id uuid.UUID) error {
	existing, err := s.repo.GetMedicineByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch existing medicine: %v", err)
	}
	if existing == nil {
		return ErrNotExist
	}

	return s.repo.DeleteMedicine(id)
}
