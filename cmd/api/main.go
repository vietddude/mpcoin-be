package main

import (
	_ "mpc/docs"
	"mpc/internal/api"
	"mpc/internal/config"
	"mpc/internal/db"
	"mpc/internal/db/redis"
	"mpc/internal/repository"
	"mpc/internal/service"
	"mpc/pkg/ethereum"
	"mpc/pkg/logger"
	"mpc/pkg/token"
	"mpc/pkg/tss"
)

// @title MPC API
// @version 1.0
// @description This is the API documentation for the MPC project.
// @host localhost:5001
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.Info("Starting application")

	// config
	logger.Info("Loading config")
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", err)
	}

	// db
	logger.Info("Initializing database")
	dbPool, err := db.InitDB(&cfg.DB)
	if err != nil {
		logger.Error("Failed to initialize database", err)
	}
	defer db.CloseDB()

	// redis
	logger.Info("Initializing Redis client")
	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Error("Failed to initialize Redis client", err)
	}
	defer redisClient.Close()

	// token
	tokenManager := token.NewTokenManager(redisClient)

	// ethereum
	ethClient, err := ethereum.NewEthClient(cfg.Eth.URL)
	if err != nil {
		logger.Error("Failed to initialize Ethereum client", err)
	}

	// tss
	tssClient, err := tss.NewTSS(redisClient)
	if err != nil {
		logger.Error("Failed to initialize TSS client", err)
	}

	// repository
	chainRepo := repository.NewChainRepository(dbPool)
	tokenRepo := repository.NewTokenRepository(dbPool)
	transactionRepo := repository.NewTransactionRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)
	walletRepo := repository.NewWalletRepository(dbPool)

	// service
	oauthClient := &service.GoogleOAuthClient{
		ClientID:     cfg.OauthClient.ClientID,
		ClientSecret: cfg.OauthClient.ClientSecret,
		RedirectURI:  cfg.OauthClient.RedirectURI,
	}
	assetService := service.NewAssetService(chainRepo, tokenRepo, redisClient)
	walletService := service.NewWalletService(walletRepo, tssClient)
	userService := service.NewUserService(userRepo, walletRepo, redisClient)
	authService := service.NewAuthService(userService, walletService, tokenManager, oauthClient)
	transactionService := service.NewTransactionService(transactionRepo, walletService, assetService, ethClient, tssClient)

	// router
	router := api.NewRouter(authService, assetService, userService, transactionService, tokenManager)

	// run router
	logger.Info("Running router")
	logger.Info("Server running on port " + cfg.Port)
	router.Run(":" + cfg.Port)
}
