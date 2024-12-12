package service

import (
	"context"
	"mpc/internal/model"
	"mpc/internal/repository"
	"mpc/pkg/errors"
	"mpc/pkg/ethereum"
	"mpc/pkg/logger"
	"strings"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
	ethClient  *ethereum.EthClient
}

func NewWalletService(walletRepo *repository.WalletRepository, ethClient *ethereum.EthClient) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		ethClient:  ethClient,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID) (model.Wallet, error) {
	// Create Ethereum wallet
	privateKeyHex, addressHex, err := s.ethClient.CreateWallet()
	if err != nil {
		logger.Error("Service:CreateWallet", err)
		return model.Wallet{}, err
	}

	addressHex = strings.ToLower(addressHex)

	// Create wallet in repository
	return s.walletRepo.CreateWallet(ctx, userID, addressHex, []byte(privateKeyHex), "Default")
}

func (s *WalletService) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (model.Wallet, error) {
	wallets, err := s.walletRepo.GetWalletsByUserID(ctx, userID)
	if err != nil {
		logger.Error("Service:GetWalletByUserID", err)
		return model.Wallet{}, err
	}
	if len(wallets) == 0 {
		return model.Wallet{}, errors.ErrWalletNotFound
	}
	return wallets[0], nil
}
