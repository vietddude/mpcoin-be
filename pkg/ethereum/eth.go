package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"mpc/pkg/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type EthClient struct {
	rpcURL string
	client *ethclient.Client
}

// NewEthClient initializes a new Ethereum client
func NewEthClient(rpcURL string) (*EthClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	return &EthClient{
		rpcURL: rpcURL,
		client: client,
	}, nil
}

// CreateWallet creates a new Ethereum wallet and returns the address and private key
func (c *EthClient) CreateWallet() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := "0x" + hex.EncodeToString(privateKeyBytes)
	addressHex := address.Hex()

	return privateKeyHex, addressHex, nil
}

// Convert amount to Wei
func (c *EthClient) convertToWei(amount string) (*big.Int, error) {
	// Parse the decimal amount
	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal amount: %w", err)
	}

	// Multiply by 10^18 for Wei conversion
	weiAmount := decimalAmount.Mul(decimal.NewFromBigInt(big.NewInt(1), 18))

	// Convert to big.Int
	wei := new(big.Int)
	wei, ok := wei.SetString(weiAmount.String(), 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert to wei")
	}

	logger.Info("amount conversion",
		zap.String("original_amount", amount),
		zap.String("decimal_amount", decimalAmount.String()),
		zap.String("wei_amount", wei.String()))

	return wei, nil
}

// CreateTransaction sends a transaction from a wallet to another address
func (c *EthClient) CreateTransaction(ctx context.Context, privateKeyHex string, to string, amount string) (string, error) {
	// Validate recipient address and amount
	toAddress, amountWei, err := c.validateTransactionInputs(to, amount)
	if err != nil {
		return "", err
	}

	// Get private key and sender address
	privateKey, fromAddress, err := c.getPrivateKeyAndSender(privateKeyHex)
	if err != nil {
		return "", err
	}

	// Fetch nonce, gas price, and chain ID
	nonce, err := c.fetchNonce(ctx, fromAddress)
	if err != nil {
		return "", err
	}
	gasPrice, err := c.fetchGasPrice(ctx)
	if err != nil {
		return "", err
	}
	chainID, err := c.fetchChainID(ctx)
	if err != nil {
		return "", err
	}

	// Build, sign, and send transaction
	txHash, err := c.buildAndSendTransaction(ctx, privateKey, nonce, toAddress, amountWei, gasPrice, chainID)
	if err != nil {
		return "", err
	}

	logger.Info("transaction sent successfully",
		zap.String("tx_hash", txHash),
		zap.String("from", fromAddress.Hex()),
		zap.String("to", to),
		zap.String("amount", amount))

	return txHash, nil
}

// validateTransactionInputs checks the recipient address and converts the amount to Wei
func (c *EthClient) validateTransactionInputs(to string, amount string) (common.Address, *big.Int, error) {
	if !common.IsHexAddress(to) {
		return common.Address{}, nil, fmt.Errorf("invalid recipient address")
	}
	toAddress := common.HexToAddress(to)

	amountWei, err := c.convertToWei(amount)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("invalid amount: %w", err)
	}
	return toAddress, amountWei, nil
}

// getPrivateKeyAndSender converts the private key and gets the sender's address
func (c *EthClient) getPrivateKeyAndSender(privateKeyHex string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	return privateKey, fromAddress, nil
}

// fetchNonce retrieves the nonce without retry
func (c *EthClient) fetchNonce(ctx context.Context, fromAddress common.Address) (uint64, error) {
	nonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch nonce: %w", err)
	}
	logger.Info("got nonce", zap.Uint64("nonce", nonce))
	return nonce, nil
}

// fetchGasPrice retrieves the gas price without retry
func (c *EthClient) fetchGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gas price: %w", err)
	}
	logger.Info("got gas price", zap.String("gas_price", gasPrice.String()))
	return gasPrice, nil
}

// fetchChainID retrieves the chain ID without retry
func (c *EthClient) fetchChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := c.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch network ID: %w", err)
	}
	return chainID, nil
}

// buildAndSendTransaction creates, signs, and sends a transaction without retry
func (c *EthClient) buildAndSendTransaction(ctx context.Context, privateKey *ecdsa.PrivateKey, nonce uint64, toAddress common.Address, amountWei *big.Int, gasPrice *big.Int, chainID *big.Int) (string, error) {
	tx := types.NewTransaction(nonce, toAddress, amountWei, 21000, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}
