package service

import (
	"context"
	"fmt"
	"mpc/internal/db/redis"
	"mpc/internal/model"
	"mpc/internal/repository"
	"mpc/pkg/cache"
	"mpc/pkg/errors"
	"mpc/pkg/utils"

	"github.com/google/uuid"
)

type AssetService struct {
	chainRepo *repository.ChainRepository
	tokenRepo *repository.TokenRepository
	cache     *cache.Cache
}

func NewAssetService(chainRepo *repository.ChainRepository, tokenRepo *repository.TokenRepository, redisClient *redis.Client) *AssetService {
	return &AssetService{
		chainRepo: chainRepo,
		tokenRepo: tokenRepo,
		cache:     cache.NewCache(redisClient, "asset"),
	}
}

// GetChains get all chains
func (s *AssetService) GetChains(ctx context.Context) ([]model.ChainResponse, error) {
	return cache.FetchOrStore(ctx, s.cache, "chains", func() ([]model.ChainResponse, error) {
		chains, err := s.chainRepo.GetChains(ctx)
		if err != nil {
			return nil, err
		}
		if len(chains) == 0 {
			return nil, errors.ErrChainNotFound
		}
		return convertChainsToResponses(chains), nil
	})
}

// GetTokensByChainID get tokens by chain id
func (s *AssetService) GetTokensByChainID(ctx context.Context, chainID int) ([]model.TokenResponse, error) {
	return cache.FetchOrStore(ctx, s.cache, fmt.Sprintf("tokens:%d", chainID), func() ([]model.TokenResponse, error) {
		chain, err := s.GetChainByChainID(ctx, chainID)
		if err != nil {
			return nil, err
		}
		tokens, err := s.tokenRepo.GetTokensByChainID(ctx, chain.ID)
		if err != nil {
			return nil, err
		}
		if len(tokens) == 0 {
			return nil, errors.ErrTokenNotFound
		}
		return convertTokensToResponses(tokens), nil
	})
}

// GetTokenBySymbol get token by symbol
func (s *AssetService) GetTokenBySymbol(ctx context.Context, chainID int, symbol string) (model.TokenResponse, error) {
	return cache.FetchOrStore(ctx, s.cache, fmt.Sprintf("token:%d:%s", chainID, symbol), func() (model.TokenResponse, error) {
		chain, err := s.GetChainByChainID(ctx, chainID)
		if err != nil {
			return model.TokenResponse{}, err
		}

		dbToken, err := s.tokenRepo.GetTokenBySymbol(ctx, chain.ID, symbol)
		if dbToken.ID == uuid.Nil {
			return model.TokenResponse{}, errors.ErrTokenNotFound
		}
		if err != nil {
			return model.TokenResponse{}, err
		}
		return utils.ToTokenResponse(dbToken), nil
	})
}

// GetChainByChainID get chain by chain id
func (s *AssetService) GetChainByChainID(ctx context.Context, chainID int) (model.ChainResponse, error) {
	return cache.FetchOrStore(ctx, s.cache, fmt.Sprintf("token:%d", chainID), func() (model.ChainResponse, error) {
		dbChain, err := s.chainRepo.GetChainByChainID(ctx, chainID)

		if dbChain.ID == uuid.Nil {
			return model.ChainResponse{}, errors.ErrChainNotFound
		}
		if err != nil {
			return model.ChainResponse{}, err
		}

		return utils.ToChainResponse(dbChain), nil
	})
}

// convertChainsToResponses convert chains to responses
func convertChainsToResponses(chains []model.Chain) []model.ChainResponse {
	result := make([]model.ChainResponse, len(chains))
	for i, chain := range chains {
		result[i] = utils.ToChainResponse(chain)
	}
	return result
}

// convertTokensToResponses convert tokens to responses
func convertTokensToResponses(tokens []model.Token) []model.TokenResponse {
	result := make([]model.TokenResponse, len(tokens))
	for i, token := range tokens {
		result[i] = utils.ToTokenResponse(token)
	}
	return result
}
