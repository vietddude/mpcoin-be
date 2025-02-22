package handler

import (
	"mpc/internal/model"
	"mpc/internal/service"
	"mpc/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(),
		authService: authService,
	}
}

// Login godoc
// @Summary      Login
// @Description  Login to the system
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.LoginRequest  true  "Login request"
// @Success      200  {object}  model.Response{payload=model.AuthResponse}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := utils.ValidateBody(c, &req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, res)
}

// Signup godoc
// @Summary      Login user
// @Description  Authenticate user and return tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.SignupRequest true "Signup credentials"
// @Success      200  {object}  model.Response{payload=model.SignupResponse}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Router       /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req model.SignupRequest
	if err := utils.ValidateBody(c, &req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.authService.Signup(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, res)
}

// Refresh godoc
// @Summary      Refresh
// @Description  Refresh the token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      model.Refresh  true  "Refresh request"
// @Success      200  {object}  model.Response{payload=model.RefreshResponse}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      401  {object}  model.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.Refresh
	if err := utils.ValidateBody(c, &req); err != nil {
		c.Error(err)
		return
	}
	res, err := h.authService.Refresh(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, res)
}
