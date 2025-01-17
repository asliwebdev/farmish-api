package models

import (
	"time"

	"github.com/google/uuid"
)

type Farm struct {
	ID uuid.UUID `json:"id"`
	CreateFarmRequest
	CreatedAt time.Time `json:"created_at"`
}

type CreateFarmRequest struct {
	Name     string    `json:"name" binding:"required"`
	Location string    `json:"location" binding:"required"`
	OwnerID  uuid.UUID `json:"owner_id" binding:"required,uuid"`
}

type CreateFarmResponse struct {
	MessageResp
	Farm `json:"farm"`
}

type UpdateFarmRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
	CreateFarmRequest
}

type UpdateFarmResp struct {
	MessageResp
	UpdateFarmRequest `json:"farm"`
}
