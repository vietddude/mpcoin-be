package handler

import (
	"mpc/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BaseHandler struct {
}

func NewBaseHandler() BaseHandler {
	return BaseHandler{}
}

// GetUserID extracts user ID from context with error handling
func (h *BaseHandler) GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.ErrUnauthorized
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.NewAppError("INVALID_USER_ID", "Invalid user ID format", http.StatusBadRequest)
	}

	return id, nil
}

// SuccessResponse send success response
func (h *BaseHandler) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"payload": data,
	})
}
