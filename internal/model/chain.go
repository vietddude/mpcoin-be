package model

import (
	"time"

	"github.com/google/uuid"
)

type Chain struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	ChainID        int       `json:"chain_id"`
	RPCURL         string    `json:"rpc_url"`
	ExplorerURL    string    `json:"explorer_url"`
	NativeCurrency string    `json:"native_currency"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ChainResponse struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name           string    `json:"name" example:"Sepolia"`
	ChainID        int       `json:"chain_id" example:"11155111"`
	RPCURL         string    `json:"rpc_url" example:"https://ethereum-sepolia-rpc.publicnode.com"`
	ExplorerURL    string    `json:"explorer_url" example:"https://sepolia.etherscan.io"`
	NativeCurrency string    `json:"native_currency" example:"ETH"`
}
