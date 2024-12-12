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

type ChainRepository struct {
	queries *db.Queries
}

func NewChainRepository(pool *pgxpool.Pool) *ChainRepository {
	return &ChainRepository{queries: db.New(pool)}
}

// GetChainByID retrieves a chain by its ID
func (r *ChainRepository) GetChainByID(ctx context.Context, id uuid.UUID) (model.Chain, error) {
	chain, err := r.queries.GetChainByID(ctx, utils.ToPgUUID(id))
	if err != nil {
		// Wrapping the error with more context
		return model.Chain{}, fmt.Errorf("failed to get chain by ID: %w", err)
	}
	return toChainModel(chain), nil
}

// GetChainByChainID retrieves a chain by its ChainID
func (r *ChainRepository) GetChainByChainID(ctx context.Context, chainID int) (model.Chain, error) {
	chain, err := r.queries.GetChainByChainID(ctx, int32(chainID))
	if err != nil {
		// Wrapping the error with more context
		return model.Chain{}, fmt.Errorf("failed to get chain by ChainID: %w", err)
	}
	return toChainModel(chain), nil
}

// GetChains retrieves all chains
func (r *ChainRepository) GetChains(ctx context.Context) ([]model.Chain, error) {
	chains, err := r.queries.GetChains(ctx)
	if err != nil {
		// Wrapping the error with more context
		return nil, fmt.Errorf("failed to get chains: %w", err)
	}

	var result []model.Chain
	for _, chain := range chains {
		result = append(result, toChainModel(chain))
	}
	return result, nil
}

// toChainModel converts a sqlc chain to a model chain
func toChainModel(sqlcChain db.Chain) model.Chain {
	return model.Chain{
		ID:             utils.ToUUID(sqlcChain.ID),
		Name:           sqlcChain.Name,
		ChainID:        int(sqlcChain.ChainID),
		RPCURL:         sqlcChain.RpcUrl,
		ExplorerURL:    sqlcChain.ExplorerUrl.String,
		NativeCurrency: sqlcChain.NativeCurrency,
		Status:         sqlcChain.Status,
		CreatedAt:      sqlcChain.CreatedAt.Time,
		UpdatedAt:      sqlcChain.UpdatedAt.Time,
	}
}
