package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `json:"id"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	ChainID     int       `json:"chain_id"`
	TxHash      string    `json:"tx_hash"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TransactionFilter struct {
	WalletID string `form:"wallet_id"`
	TokenID  string `form:"token_id"`
	Status   string `form:"status"`
	FromDate string `form:"from_date"`
	ToDate   string `form:"to_date"`
}

type Pagination struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

type TransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	PageSize     int           `json:"page_size"`
	TotalPages   int           `json:"total_pages"`
}

type CreateAndSubmitTransactionRequest struct {
	FromAddress string `json:"from_address" validate:"required"`
	ToAddress   string `json:"to_address" validate:"required"`
	ChainID     int    `json:"chain_id" validate:"required"`
	Symbol      string `json:"symbol" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	ShareData   string `json:"share_data" validate:"required"`
}
