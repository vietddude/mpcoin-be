package service

import (
	"context"
	"mpc/internal/model"
	"mpc/internal/repository"
	"mpc/pkg/errors"
	"mpc/pkg/logger"
	"mpc/pkg/tss"
	"strings"

	"github.com/google/uuid"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
	tssClient  *tss.TSS
}

func NewWalletService(walletRepo *repository.WalletRepository, tssClient *tss.TSS) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		tssClient:  tssClient,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID) (model.Wallet, []byte, error) {
	// Create Ethereum wallet
	shareData, addressHex, err := s.tssClient.CreateWallet(ctx, userID.String())
	if err != nil {
		logger.Error("Service:CreateWallet", err)
		return model.Wallet{}, nil, err
	}
	addressHex = strings.ToLower(addressHex)

	// Create wallet in repository
	wallet, err := s.walletRepo.CreateWallet(ctx, userID, addressHex, []byte(""), "Default")
	if err != nil {
		logger.Error("Service:CreateWallet", err)
		return model.Wallet{}, nil, err
	}
	return wallet, shareData, nil
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
