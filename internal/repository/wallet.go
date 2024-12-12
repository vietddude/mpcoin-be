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

type WalletRepository struct {
	queries *db.Queries
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{queries: db.New(pool)}
}

// CreateWallet creates a new wallet
func (r *WalletRepository) CreateWallet(ctx context.Context, userID uuid.UUID, address string, encryptedPrivateKey []byte, name string) (model.Wallet, error) {
	wallet, err := r.queries.CreateWallet(ctx, db.CreateWalletParams{
		UserID:              utils.ToPgUUID(userID),
		Address:             address,
		EncryptedPrivateKey: encryptedPrivateKey,
		Name:                utils.ToPgText(name),
		Status:              "active",
		CreatedAt:           utils.CurrentPgTimestamp(),
		UpdatedAt:           utils.CurrentPgTimestamp(),
	})
	if err != nil {
		return model.Wallet{}, fmt.Errorf("failed to create wallet: %w", err)
	}
	return toWalletModel(wallet), nil
}

// GetWalletByID retrieves a wallet by its ID
func (r *WalletRepository) GetWalletByID(ctx context.Context, id uuid.UUID) (model.Wallet, error) {
	wallet, err := r.queries.GetWalletByID(ctx, utils.ToPgUUID(id))
	if err != nil {
		return model.Wallet{}, fmt.Errorf("failed to get wallet by ID: %w", err)
	}
	return toWalletModel(wallet), nil
}

// GetWalletsByUserID retrieves wallets by user ID
func (r *WalletRepository) GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wallet, error) {
	wallets, err := r.queries.GetWalletsByUserID(ctx, utils.ToPgUUID(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets by user ID: %w", err)
	}
	var result []model.Wallet
	for _, wallet := range wallets {
		result = append(result, toWalletModel(wallet))
	}
	return result, nil
}

// GetWalletByAddress retrieves a wallet by its address
func (r *WalletRepository) GetWalletByAddress(ctx context.Context, address string) (model.Wallet, error) {
	wallet, err := r.queries.GetWalletByAddress(ctx, address)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("failed to get wallet by address: %w", err)
	}
	return toWalletModel(wallet), nil
}

// toWalletModel converts a sqlc wallet to a model wallet
func toWalletModel(sqlcWallet db.Wallet) model.Wallet {
	return model.Wallet{
		ID:                  utils.ToUUID(sqlcWallet.ID),
		UserID:              utils.ToUUID(sqlcWallet.UserID),
		Address:             sqlcWallet.Address,
		EncryptedPrivateKey: string(sqlcWallet.EncryptedPrivateKey),
		Name:                utils.ToText(sqlcWallet.Name),
		Status:              sqlcWallet.Status,
		CreatedAt:           sqlcWallet.CreatedAt.Time,
		UpdatedAt:           sqlcWallet.UpdatedAt.Time,
	}
}
