package models

import "github.com/google/uuid"

type Medicine struct {
	MedicineWithoutTime
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MedicineWithoutTime struct {
	ID uuid.UUID `json:"id"`
	MedicineReq
}

type MedicineReq struct {
	FarmID        uuid.UUID `json:"farm_id" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	SuitableFor   []string  `json:"suitable_for" binding:"required"`
	UnitOfMeasure string    `json:"unit_of_measure" binding:"required"`
	Quantity      float64   `json:"quantity" binding:"required,gt=0"`
	MinThreshold  float64   `json:"min_threshold" binding:"required,gt=0"`
}

type MedicineResp struct {
	MessageResp
	MedicineWithoutTime `json:"medicine"`
}
