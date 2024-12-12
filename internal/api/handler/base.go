package handler

import (
	"mpc/pkg/errors"
	"mpc/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BaseHandler struct {
}

func NewBaseHandler() BaseHandler {
	return BaseHandler{}
}

func (h *BaseHandler) GetUserID(c *gin.Context) uuid.UUID {
	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":      "Unauthorized",
			"error_code": "UNAUTHORIZED",
		})
	}
	return userID.(uuid.UUID)
}

// HandleError handle error and send error response
func (h *BaseHandler) HandleError(c *gin.Context, err error) {
	logger.Error("error", err)
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

// SuccessResponse send success response
func (h *BaseHandler) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"payload": data,
	})
}
