package middleware

import (
	"mpc/pkg/logger"
	"mpc/pkg/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verify JWT token from header "Authorization"
func AuthMiddleware(jwtService *token.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header Authorization
		jwtToken := c.GetHeader("Authorization")
		if jwtToken == "" {
			logger.Info("AuthMiddleware: jwtToken is empty")
			raiseUnauthorizedError(c)
			return
		}

		// Remove "Bearer " if it exists
		jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

		// Verify token
		userID, err := jwtService.VerifyToken(c.Request.Context(), jwtToken, token.TokenTypeAccess)
		if err != nil {
			logger.Error("AuthMiddleware: Invalid or expired token", err)
			raiseUnauthorizedError(c)
			return
		}
		// Save user_id to context to be used by subsequent handlers
		c.Set("user_id", userID)

		// Continue processing request
		c.Next()
	}
}

func raiseUnauthorizedError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error":      "Unauthorized",
		"error_code": "UNAUTHORIZED",
	})
}
