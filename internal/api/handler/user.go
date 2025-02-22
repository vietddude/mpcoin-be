package handler

import (
	"mpc/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	BaseHandler
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(),
		userService: userService,
	}
}

// GetUser godoc
// @Summary      Get user
// @Description  Get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Response{payload=model.UserResponse}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Router       /user [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := h.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}
	// Handle the request using the optimized request handler
	res, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, res)
}
