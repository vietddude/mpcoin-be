package repository

import (
	"context"
	"fmt"
	db "mpc/internal/db/sqlc"
	"mpc/internal/model"
	"mpc/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	queries *db.Queries
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{queries: db.New(pool)}
}

// CreateTransaction creates a new transaction
func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	tx, err := r.queries.CreateTransaction(ctx, db.CreateTransactionParams{
		ChainID:     int32(transaction.ChainID),
		FromAddress: transaction.FromAddress,
		ToAddress:   transaction.ToAddress,
		TxHash:      transaction.TxHash,
		CreatedAt:   utils.CurrentPgTimestamp(),
		UpdatedAt:   utils.CurrentPgTimestamp(),
	})
	if err != nil {
		return model.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}
	return toTransactionModel(tx), nil
}

// GetTransactionsByWalletAddress retrieves transactions by wallet ID
func (r *TransactionRepository) GetTransactionsByWalletAddress(ctx context.Context, walletAddress string, chainID int, limit int, offset int) ([]model.Transaction, error) {
	transactions, err := r.queries.GetTransactionsByWalletAddress(ctx, db.GetTransactionsByWalletAddressParams{
		FromAddress: walletAddress,
		Column2:     int32(chainID),
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by wallet address: %w", err)
	}

	var result []model.Transaction
	for _, tx := range transactions {
		result = append(result, toTransactionModel(tx))
	}
	return result, nil
}

// GetTransactionByID retrieves a transaction by its ID
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (model.Transaction, error) {
	transaction, err := r.queries.GetTransactionByID(ctx, utils.ToPgUUID(id))
	if err != nil {
		return model.Transaction{}, fmt.Errorf("failed to get transaction by ID: %w", err)
	}
	return toTransactionModel(transaction), nil
}

// GetTransactionCount retrieves the number of transactions for a given wallet ID and chain ID
func (r *TransactionRepository) GetTransactionCount(ctx context.Context, walletAddress string, chainID int) (int, error) {
	count, err := r.queries.GetTransactionCount(ctx, db.GetTransactionCountParams{
		FromAddress: walletAddress,
		Column2:     int32(chainID),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}
	return int(count), nil
}

// toTransactionModel converts a sqlc transaction to a model transaction
func toTransactionModel(sqlcTransaction db.Transaction) model.Transaction {
	return model.Transaction{
		ID:          utils.ToUUID(sqlcTransaction.ID),
		ChainID:     int(sqlcTransaction.ChainID),
		FromAddress: sqlcTransaction.FromAddress,
		ToAddress:   sqlcTransaction.ToAddress,
		TxHash:      sqlcTransaction.TxHash,
		CreatedAt:   sqlcTransaction.CreatedAt.Time,
		UpdatedAt:   sqlcTransaction.UpdatedAt.Time,
	}
}
