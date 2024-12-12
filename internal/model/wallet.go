package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"user_id"`
	Address             string    `json:"address"`
	EncryptedPrivateKey string    `json:"encrypted_private_key"`
	Name                string    `json:"name"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type WalletResponse struct {
	ID      uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID  uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Address string    `json:"address" example:"0x0000000000000000000000000000000000000000"`
	Name    string    `json:"name" example:"My Wallet"`
}
