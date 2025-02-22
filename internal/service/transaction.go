package service

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"mpc/internal/model"
	"mpc/internal/repository"
	"mpc/pkg/errors"
	"mpc/pkg/ethereum"
	"mpc/pkg/logger"
	"mpc/pkg/tss"
	"mpc/pkg/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

type TransactionService struct {
	txnRepo       *repository.TransactionRepository
	assetService  *AssetService
	walletService *WalletService
	ethClient     *ethereum.EthClient
	tssClient     *tss.TSS
}

func NewTransactionService(
	txnRepo *repository.TransactionRepository,
	walletService *WalletService,
	assetService *AssetService,
	ethClient *ethereum.EthClient,
	tssClient *tss.TSS,
) *TransactionService {
	return &TransactionService{
		txnRepo:       txnRepo,
		assetService:  assetService,
		walletService: walletService,
		ethClient:     ethClient,
		tssClient:     tssClient,
	}
}

// GetTransactions retrieves a list of transactions based on the provided parameters.
func (s *TransactionService) GetTransactions(
	ctx context.Context,
	userID uuid.UUID,
	chainID int,
	walletAddress string,
	page, pageSize int,
) (model.TransactionListResponse, error) {
	// Validate pagination
	page, pageSize, err := utils.ValidatePagination(page, pageSize)
	if err != nil {
		return model.TransactionListResponse{}, errors.ErrInvalidRequest
	}

	walletAddress = strings.ToLower(walletAddress)

	// Validate wallet and chain
	if _, err = s.walletService.walletRepo.GetWalletByAddress(ctx, walletAddress); err != nil {
		return model.TransactionListResponse{}, errors.ErrWalletNotFound
	}

	if _, err = s.assetService.chainRepo.GetChainByChainID(ctx, chainID); err != nil {
		return model.TransactionListResponse{}, errors.ErrInvalidChainID
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Fetch transactions
	transactions, err := s.txnRepo.GetTransactionsByWalletAddress(ctx, walletAddress, chainID, pageSize, offset)
	if err != nil {
		return model.TransactionListResponse{}, errors.ErrTransactionNotFound
	}

	// Get total transaction count
	total, err := s.txnRepo.GetTransactionCount(ctx, walletAddress, chainID)
	if err != nil {
		return model.TransactionListResponse{}, errors.ErrTransactionNotFound
	}

	// Calculate total pages
	totalPages := (total + pageSize - 1) / pageSize

	return model.TransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}, nil
}

// CreateAndSubmitTransaction creates and submits a transaction.
func (s *TransactionService) CreateAndSubmitTransaction(
	ctx context.Context,
	userID uuid.UUID,
	req model.CreateAndSubmitTransactionRequest,
) (model.Transaction, error) {
	// Validate request
	if err := s.validateRequest(req); err != nil {
		return model.Transaction{}, err
	}

	req.FromAddress = strings.ToLower(req.FromAddress)
	req.ToAddress = strings.ToLower(req.ToAddress)

	// Fetch wallet
	wallet, err := s.walletService.GetWalletByUserID(ctx, userID)
	if err != nil {
		logger.Error("failed to get wallet by user ID", err)
		return model.Transaction{}, errors.ErrInvalidRequest
	}

	if !strings.EqualFold(wallet.Address, req.FromAddress) {
		logger.Warn("wallet address does not match the from address in the request")
		return model.Transaction{}, errors.ErrInvalidRequest
	}

	// Check if the wallet has enough balance
	enough, err := s.ethClient.IsEnoughBalance(ctx, req.FromAddress, req.Amount)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("failed to check balance: %w", err)
	}
	if !enough {
		return model.Transaction{}, errors.ErrInssuficientBalance
	}

	hash, err := s.handleTxn(ctx, userID.String(), req)
	if err != nil {
		return model.Transaction{}, err
	}
	// Create transaction record in the database
	return s.createTransactionRecord(ctx, req.FromAddress, req.ToAddress, hash, req.ChainID)
}

// validateRequest validates the transaction request.
func (s *TransactionService) validateRequest(req model.CreateAndSubmitTransactionRequest) error {
	if !common.IsHexAddress(req.FromAddress) || !common.IsHexAddress(req.ToAddress) {
		return errors.ErrInvalidAddress
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil || amount <= 0 {
		return errors.ErrInvalidAmount
	}

	return nil
}

// createTransactionRecord creates a new transaction record in the database.
func (s *TransactionService) createTransactionRecord(
	ctx context.Context,
	fromAddress, toAddress, txHash string, chainID int,
) (model.Transaction, error) {
	txn := model.Transaction{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		TxHash:      txHash,
		ChainID:     chainID,
	}

	// Save transaction in the repository
	createdTxn, err := s.txnRepo.CreateTransaction(ctx, txn)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("failed to create transaction record: %w", err)
	}

	return createdTxn, nil
}

func (s *TransactionService) handleTxn(ctx context.Context, userID string, req model.CreateAndSubmitTransactionRequest) (string, error) {
	// Validate chain ID (Sepolia testnet: 11155111)
	chainID := big.NewInt(11155111)
	if req.ChainID != 0 && req.ChainID != int(chainID.Int64()) {
		return "", fmt.Errorf("invalid chain ID: got %d, want %d", req.ChainID, chainID.Uint64())
	}

	// Validate addresses
	if !common.IsHexAddress(req.FromAddress) || !common.IsHexAddress(req.ToAddress) {
		return "", fmt.Errorf("invalid address format")
	}

	// Tạo transaction
	tx, err := s.ethClient.CreateTransaction(ctx, req.FromAddress, req.ToAddress, req.Amount)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Lấy transaction hash
	signer := types.NewEIP155Signer(chainID)
	txHash := signer.Hash(tx)

	// Ký bằng TSS (nhận chữ ký DER)
	derSig, err := s.tssClient.Sign(ctx, userID, req.ShareData, txHash.Bytes())
	if err != nil {
		return "", fmt.Errorf("TSS signing failed: %w", err)
	}

	fmt.Print("from address: ", req.FromAddress)

	sig, err := utils.ConvertDERToEthSignature(derSig, txHash.Bytes(), req.FromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert DER signature: %w", err)
	}
	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	// Gửi transaction
	txHashSent, err := s.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	return txHashSent, nil
}
