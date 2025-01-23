package models

import (
	"time"

	"github.com/google/uuid"
)

type MedicalRecordWithoutTime struct {
	ID uuid.UUID `json:"id"`
	MedicalRecordReq
}

type MedicalRecordReq struct {
	AnimalID   uuid.UUID `json:"animal_id" binding:"required,uuid"`
	MedicineID uuid.UUID `json:"medicine_id" binding:"required,uuid"`
	UpdateMedicalRecordReq
}

type MedicalRecordResp struct {
	MessageResp
	MedicalRecordWithoutTime `json:"medical_record"`
}

type UpdateMedicalRecordReq struct {
	Quantity      float64   `json:"quantity" binding:"required,gt=0"`
	TreatmentDate time.Time `json:"treatment_date" binding:"required"`
	Notes         string    `json:"notes" binding:"max=500"`
}

type MedicalRecordDetailed struct {
	ID            string         `json:"id"`
	Animal        AnimalDetail   `json:"animal"`
	Medicine      MedicineDetail `json:"medicine"`
	Quantity      float64        `json:"quantity"`
	TreatmentDate time.Time      `json:"treatment_date"`
	Notes         string         `json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
}

type MedicineDetail struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	SuitableFor   []string  `json:"suitable_for"`
	UnitOfMeasure string    `json:"unit_of_measure"`
}
