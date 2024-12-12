package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TokenTypeNative = "NATIVE"
	TokenTypeERC20  = "ERC20"
)

type Token struct {
	ID              uuid.UUID `json:"id"`
	ChainID         uuid.UUID `json:"chain_id"`
	ContractAddress string    `json:"contract_address"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	Decimals        int32     `json:"decimals"`
	LogoURL         string    `json:"logo_url"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TokenResponse struct {
	ID              uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ChainID         uuid.UUID `json:"chain_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ContractAddress string    `json:"contract_address" example:"0x0000000000000000000000000000000000000000"`
	Name            string    `json:"name" example:"Ethereum"`
	Symbol          string    `json:"symbol" example:"ETH"`
	Decimals        int32     `json:"decimals" example:"18"`
	LogoURL         string    `json:"logo_url" example:"https://example.com/logo.png"`
	Type            string    `json:"type" example:"ERC20"`
}
