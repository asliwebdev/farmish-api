package services

import (
	"errors"
	"farmish/internal/models"
	"farmish/internal/repository"

	"github.com/google/uuid"
)

type MedicalRecordService struct {
	medicalRecordRepo *repository.MedicalRecordRepository
	animalRepo        *repository.AnimalRepository
	medicineRepo      *repository.MedicineRepository
}

func NewMedicalRecordService(medicalRecordRepo *repository.MedicalRecordRepository,
	animalRepo *repository.AnimalRepository,
	medicineRepo *repository.MedicineRepository) *MedicalRecordService {
	return &MedicalRecordService{
		medicalRecordRepo: medicalRecordRepo,
		animalRepo:        animalRepo,
		medicineRepo:      medicineRepo,
	}
}

var (
	ErrMedicalRecordNotFound = errors.New("medical record not found")
)

func (s *MedicalRecordService) CreateMedicalRecord(record *models.MedicalRecordWithoutTime) error {
	animal, err := s.animalRepo.GetAnimalByID(record.AnimalID)
	if err != nil {
		return err
	} else if animal == nil {
		return ErrAnimalNotFound
	}

	medicine, err := s.medicineRepo.GetMedicineByID(record.MedicineID)
	if err != nil {
		return err
	} else if medicine == nil {
		return ErrMedicineNotExist
	}

	if medicine.Quantity < record.Quantity {
		return ErrInsufficientQuantity
	}

	newQuantity := medicine.Quantity - record.Quantity

	record.ID = uuid.New()

	return s.medicalRecordRepo.CreateMedicalRecord(record, newQuantity)
}

func (s *MedicalRecordService) GetMedicalRecordByID(recordID uuid.UUID) (*models.MedicalRecordDetailed, error) {
	return s.medicalRecordRepo.GetMedicalRecordByID(recordID)
}

func (s *MedicalRecordService) GetMedicalRecordsByAnimalID(animalID uuid.UUID) ([]*models.MedicalRecordDetailed, error) {
	return s.medicalRecordRepo.GetMedicalRecordsByAnimalID(animalID)
}

func (s *MedicalRecordService) UpdateMedicalRecord(record *models.MedicalRecordWithoutTime) error {
	return s.medicalRecordRepo.UpdateMedicalRecord(record)
}

func (s *MedicalRecordService) DeleteMedicalRecord(recordID uuid.UUID) error {
	return s.medicalRecordRepo.DeleteMedicalRecord(recordID)
}
