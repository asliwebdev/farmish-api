package models

import (
	"time"

	"github.com/google/uuid"
)

type FeedingRecordReq struct {
	AnimalID uuid.UUID `json:"animal_id" binding:"required,uuid"`
	FoodID   uuid.UUID `json:"food_id" binding:"required,uuid"`
	UpdateFeedRecordReq
}

type FeedingRecordResp struct {
	MessageResp
	FeedingRecordWithoutTime `json:"feeding_record"`
}

type FeedingRecordWithoutTime struct {
	ID uuid.UUID `json:"id"`
	FeedingRecordReq
}

type UpdateFeedRecordReq struct {
	Quantity float64   `json:"quantity" binding:"required,gt=0"`
	FedAt    time.Time `json:"fed_at" binding:"required"`
	Notes    string    `json:"notes" binding:"max=500"`
}

type FeedingRecordDetailed struct {
	FeedingRecordID uuid.UUID    `json:"feeding_record_id"`
	Quantity        float64      `json:"quantity"`
	FedAt           time.Time    `json:"fed_at"`
	Notes           string       `json:"notes,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	Animal          AnimalDetail `json:"animal"`
	Food            FoodDetail   `json:"food"`
}

type AnimalDetail struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Weight       float64   `json:"weight"`
	HealthStatus string    `json:"health_status"`
}

type FoodDetail struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	SuitableFor   []string  `json:"suitable_for"`
	UnitOfMeasure string    `json:"unit_of_measure"`
}
