package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	SignUpRequest
	CreatedAt time.Time `json:"created_at"`
}

type SignUpRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,min=9"`
	LoginRequest
}

type SignUpResponse struct {
	MessageResp
	LoginResponse
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	ID    uuid.UUID `json:"user_id"`
	Token string    `json:"token"`
}

type UpdateUserSwag struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=9"`
	Password    string `json:"password"`
}

type UpdateUser struct {
	ID uuid.UUID `json:"id"`
	UpdateUserSwag
}

type UpdateUserResp struct {
	MessageResp
	UpdateUser `json:"user"`
}
