package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"mpc/pkg/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// CreateTransaction sends a transaction from a wallet to another address
func (c *EthClient) CreateTransaction(ctx context.Context, fromAddressHex string, to string, amount string) (*types.Transaction, error) {
	// Validate recipient address and amount
	toAddress, amountWei, err := c.validateTransactionInputs(to, amount)
	if err != nil {
		return nil, err
	}

	fromAddress := common.HexToAddress(fromAddressHex)

	// Fetch nonce, gas price, and chain ID
	nonce, err := c.fetchNonce(ctx, fromAddress)
	if err != nil {
		return nil, err
	}
	gasPrice, err := c.fetchGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(nonce, toAddress, amountWei, 210000, gasPrice, nil)

	return tx, nil
}

func (c *EthClient) SendTransaction(ctx context.Context, signedTx *types.Transaction) (string, error) {
	err := c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	return signedTx.Hash().Hex(), nil
}

func (c *EthClient) IsEnoughBalance(ctx context.Context, address string, amount string) (bool, error) {
	// Validate recipient address and amount
	_, amountWei, err := c.validateTransactionInputs(address, amount)
	if err != nil {
		return false, err
	}

	balance, err := c.client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return false, fmt.Errorf("failed to fetch balance: %w", err)
	}

	return balance.Cmp(amountWei) >= 0, nil
}

// validateTransactionInputs checks the recipient address and converts the amount to Wei
func (c *EthClient) validateTransactionInputs(address string, amount string) (common.Address, *big.Int, error) {
	if !common.IsHexAddress(address) {
		return common.Address{}, nil, fmt.Errorf("invalid recipient address")
	}
	addr := common.HexToAddress(address)

	amountWei, err := c.convertToWei(amount)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("invalid amount: %w", err)
	}
	return addr, amountWei, nil
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
