package utils

import (
	"mpc/internal/model"

	"github.com/google/uuid"
)

func ToChainResponse(chain model.Chain) model.ChainResponse {
	return model.ChainResponse{
		ID:             chain.ID,
		Name:           chain.Name,
		ChainID:        chain.ChainID,
		RPCURL:         chain.RPCURL,
		ExplorerURL:    chain.ExplorerURL,
		NativeCurrency: chain.NativeCurrency,
	}
}

func ToTokenResponse(token model.Token) model.TokenResponse {
	return model.TokenResponse{
		ID:              token.ID,
		ChainID:         token.ChainID,
		ContractAddress: token.ContractAddress,
		Name:            token.Name,
		Symbol:          token.Symbol,
		Decimals:        token.Decimals,
		LogoURL:         token.LogoURL,
		Type:            token.Type,
	}
}

// Response constructors with validation
func ToUserResponse(user model.User) model.UserResponse {
	if user.ID == uuid.Nil {
		return model.UserResponse{}
	}
	return model.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func ToWalletResponse(wallet model.Wallet) model.WalletResponse {
	if wallet.ID == uuid.Nil {
		return model.WalletResponse{}
	}
	return model.WalletResponse{
		ID:      wallet.ID,
		UserID:  wallet.UserID,
		Address: wallet.Address,
	}
}
