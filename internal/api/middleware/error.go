package middleware

import (
	"mpc/pkg/errors"
	"mpc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// Log the error
			logger.Error("error", err)

			// Check if the error is an AppError
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.Status, gin.H{
					"error":      appErr.Message,
					"error_code": appErr.Code,
				})
				return
			}

			// Default error handling for non-AppError types
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":      "Internal server error",
				"error_code": "INTERNAL_SERVER_ERROR",
			})
		}
	}
}
