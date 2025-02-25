package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"mpc/internal/config"
	"mpc/internal/db"
	"mpc/internal/db/redis"
	"mpc/internal/model"
	"mpc/internal/repository"
	"mpc/pkg/logger"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL = "https://sepolia.infura.io/v3/6c89fb7fa351451f939eea9da6bee755"
)

var (
	ctx                = context.Background()
	redisClient        *redis.Client
	monitoredAddresses map[common.Address]bool
	txnRepo            *repository.TransactionRepository
	walletRepo         *repository.WalletRepository
)

func main() {
	logger.Info("Starting worker")

	// Load config
	logger.Info("Loading config")
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", err)
	}

	// Initialize database
	logger.Info("Initializing database")
	dbPool, err := db.InitDB(&cfg.DB)
	if err != nil {
		logger.Error("Failed to initialize database", err)
	}
	defer db.CloseDB()

	txnRepo = repository.NewTransactionRepository(dbPool)
	walletRepo = repository.NewWalletRepository(dbPool)

	// Initialize Redis
	logger.Info("Initializing Redis client")
	redisClient, err = redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Error("Failed to initialize Redis client", err)
	}
	defer redisClient.Close()

	// Load addresses into Redis
	if err := loadAddressesToRedis(); err != nil {
		log.Fatalf("Failed to load addresses: %v", err)
	}

	go updateCachePeriodically()

	// Connect to Ethereum client
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to Infura: %v", err)
	}
	defer client.Close()

	fmt.Println("Starting transaction scanner...")

	for {
		checkLatestBlock(client)
		time.Sleep(10 * time.Second)
	}
}

func loadAddressesToRedis() error {
	rows, err := walletRepo.GetAllAddresses(ctx)
	if err != nil {
		return err
	}

	var addresses []string
	addresses = append(addresses, rows...)
	fmt.Printf("Loaded %d addresses\n", len(addresses))
	if len(addresses) > 0 {
		redisClient.Del(ctx, "monitored_addresses")
		redisClient.SAdd(ctx, "monitored_addresses", addresses)
	}
	return nil
}

func getMonitoredAddressesFromRedis() (map[common.Address]bool, error) {
	addresses, err := redisClient.SMembers(ctx, "monitored_addresses").Result()
	if err != nil {
		return nil, err
	}

	addressMap := make(map[common.Address]bool)
	for _, addr := range addresses {
		addressMap[common.HexToAddress(addr)] = true
	}
	return addressMap, nil
}

func updateCachePeriodically() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := loadAddressesToRedis(); err != nil {
			log.Printf("Error updating Redis cache: %v", err)
		} else {
			log.Println("Updated monitored addresses in Redis")
		}
	}
}

func checkLatestBlock(client *ethclient.Client) {
	monitoredAddresses, _ = getMonitoredAddressesFromRedis()

	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Printf("Error getting latest block: %v", err)
		return
	}

	fmt.Printf("Scanning Block #%d...\n", block.NumberU64())

	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}

		from, err := client.TransactionSender(ctx, tx, block.Header().Hash(), 0)
		if err != nil {
			log.Printf("Error getting sender: %v", err)
			continue
		}

		to := *tx.To()
		if monitoredAddresses[from] || monitoredAddresses[to] {
			fmt.Printf("Transaction Found! Hash: %s, From: %s, To: %s, Value: %s ETH\n",
				tx.Hash().Hex(), from.Hex(), to.Hex(), weiToEth(tx.Value()))

			// Save transaction to database
			txn := model.Transaction{
				TxHash:      strings.ToLower(tx.Hash().Hex()),
				FromAddress: strings.ToLower(from.Hex()),
				ToAddress:   strings.ToLower(to.Hex()),
				ChainID:     11155111,
			}
			if _, err := txnRepo.CreateTransaction(ctx, txn); err != nil {
				log.Printf("Error saving transaction: %v", err)
			}
		}
	}
}

func weiToEth(wei *big.Int) string {
	ethValue := new(big.Float).SetInt(wei)
	ethValue.Quo(ethValue, big.NewFloat(1e18))
	return ethValue.Text('f', 6)
}
