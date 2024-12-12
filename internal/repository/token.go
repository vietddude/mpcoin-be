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

type TokenRepository struct {
	queries *db.Queries
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{queries: db.New(pool)}
}

// GetTokensByChainID retrieves all tokens for a given chain ID
func (r *TokenRepository) GetTokensByChainID(ctx context.Context, chainID uuid.UUID) ([]model.Token, error) {
	tokens, err := r.queries.GetTokensByChain(ctx, utils.ToPgUUID(chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens by chain ID: %w", err)
	}
	var result []model.Token
	for _, token := range tokens {
		result = append(result, toTokenModel(token))
	}
	return result, nil
}

// GetTokenByContractAddress retrieves a token by its contract address
func (r *TokenRepository) GetTokenByContractAddress(ctx context.Context, contractAddress string) (model.Token, error) {
	token, err := r.queries.GetTokenByContractAddress(ctx, contractAddress)
	if err != nil {
		return model.Token{}, fmt.Errorf("failed to get token by contract address: %w", err)
	}
	return toTokenModel(token), nil
}

// GetTokenBySymbol retrieves a token by its symbol
func (r *TokenRepository) GetTokenBySymbol(ctx context.Context, chainID uuid.UUID, symbol string) (model.Token, error) {
	token, err := r.queries.GetTokenBySymbol(ctx, db.GetTokenBySymbolParams{
		ChainID: utils.ToPgUUID(chainID),
		Symbol:  symbol,
	})
	if err != nil {
		return model.Token{}, fmt.Errorf("failed to get token by symbol: %w", err)
	}
	return toTokenModel(token), nil
}

// GetTokenByID retrieves a token by its ID
func (r *TokenRepository) GetTokenByID(ctx context.Context, id uuid.UUID) (model.Token, error) {
	token, err := r.queries.GetTokenByID(ctx, utils.ToPgUUID(id))
	if err != nil {
		return model.Token{}, fmt.Errorf("failed to get token by ID: %w", err)
	}
	return toTokenModel(token), nil
}

// toTokenModel converts a sqlc token to a model token
func toTokenModel(sqlcToken db.Token) model.Token {
	return model.Token{
		ID:              utils.ToUUID(sqlcToken.ID),
		ChainID:         utils.ToUUID(sqlcToken.ChainID),
		ContractAddress: sqlcToken.ContractAddress,
		Name:            sqlcToken.Name,
		Symbol:          sqlcToken.Symbol,
		Decimals:        sqlcToken.Decimals,
		LogoURL:         sqlcToken.LogoUrl.String,
		Type:            sqlcToken.Type,
		Status:          sqlcToken.Status,
		CreatedAt:       sqlcToken.CreatedAt.Time,
		UpdatedAt:       sqlcToken.UpdatedAt.Time,
	}
}
