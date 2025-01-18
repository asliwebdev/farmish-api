package models

import (
	"time"

	"github.com/google/uuid"
)

type Animal struct {
	ID uuid.UUID `json:"id"`
	CreateAnimalReq
	LastFed     time.Time `json:"last_fed"`
	LastWatered time.Time `json:"last_watered"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AnimalWithoutTime struct {
	ID uuid.UUID `json:"id"`
	CreateAnimalReq
}

type CreateAnimalReq struct {
	FarmID       uuid.UUID `json:"farm_id" binding:"required,uuid"`
	Name         string    `json:"name"`
	Type         string    `json:"type" binding:"required"`
	Weight       float64   `json:"weight" binding:"required"`
	HealthStatus string    `json:"health_status"`
	DateOfBirth  time.Time `json:"date_of_birth"`
}

type CreateAnimalResp struct {
	MessageResp
	AnimalWithoutTime `json:"animal"`
}

type UpdateAnimalReq struct {
	ID           uuid.UUID `json:"id" binding:"required,uuid"`
	Name         string    `json:"name"`
	Type         string    `json:"type" binding:"required"`
	Weight       float64   `json:"weight" binding:"required"`
	HealthStatus string    `json:"health_status" binding:"required"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	LastFed      time.Time `json:"last_fed" binding:"required"`
	LastWatered  time.Time `json:"last_watered" binding:"required"`
}

type UpdateAnimalResp struct {
	MessageResp
	UpdateAnimalReq `json:"animal"`
}
