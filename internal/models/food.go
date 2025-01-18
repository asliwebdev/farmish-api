package models

import "github.com/google/uuid"

type Food struct {
	FoodWithoutTime
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type FoodWithoutTime struct {
	ID uuid.UUID `json:"id"`
	AddFoodReq
}

type AddFoodReq struct {
	FarmID        uuid.UUID `json:"farm_id" binding:"required,uuid"`
	Name          string    `json:"name" binding:"required"`
	SuitableFor   []string  `json:"suitable_for" binding:"required"`
	UnitOfMeasure string    `json:"unit_of_measure" binding:"required"`
	Quantity      float64   `json:"quantity" binding:"required,gt=0"`
	MinThreshold  float64   `json:"min_threshold" binding:"required,gt=0"`
}

type AddFoodResp struct {
	MessageResp
	FoodWithoutTime `json:"food"`
}

type UpdateFoodReq struct {
	ID uuid.UUID `json:"id" binding:"required,uuid"`
	AddFoodReq
}

type UpdateFoodResp struct {
	MessageResp
	UpdateFoodReq `json:"food"`
}
