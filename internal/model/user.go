package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email string    `json:"email" example:"user@example.com"`
}

type GetUserResponse struct {
	User   UserResponse   `json:"user"`
	Wallet WalletResponse `json:"wallet"`
}
