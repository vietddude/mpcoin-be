package api

import (
	"mpc/internal/api/handler"
	"mpc/internal/api/middleware"
	"mpc/internal/service"
	"mpc/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	authService *service.AuthService,
	assetService *service.AssetService,
	userService *service.UserService,
	txnService *service.TransactionService,
	tokenManager *token.TokenManager,
) *gin.Engine {
	// Disable default logger
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	// Apply CORS middlewareelopment
	// router.Use(middleware.Cors())

	// Use our custom logger and recovery middleware
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	assetHandler := handler.NewAssetHandler(assetService)
	userHandler := handler.NewUserHandler(userService)
	txnHandler := handler.NewTransactionHandler(txnService)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/health", healthHandler.HealthCheck)

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/refresh", authHandler.Refresh)
		}

		assets := v1.Group("/assets")
		{
			assets.GET("/chains", assetHandler.GetChains)
			assets.GET("/chains/:chain_id/tokens", assetHandler.GetTokensByChainID)
			assets.GET("/chains/:chain_id/tokens/:symbol", assetHandler.GetTokenBySymbol)
		}

		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(tokenManager))
		{
			users.GET("/me", userHandler.GetUser)
		}

		transactions := v1.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(tokenManager))
		{
			transactions.GET("", txnHandler.GetTransactions)
			transactions.POST("/", txnHandler.CreateAndSubmitTransaction)
		}

		// Redirect to swagger docs
		v1.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/api/v1/swagger/index.html")
		})

	}

	// Redirect to swagger docs
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/v1/swagger/index.html")
	})
	return router
}
