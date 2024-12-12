package handler

import "github.com/gin-gonic/gin"

type HealthHandler struct {
	BaseHandler
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(),
	}
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the server is running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Response{payload=map[string]string}
// @Router       /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	h.SuccessResponse(c, gin.H{"message": "OK"})
}
