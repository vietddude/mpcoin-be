package handler

import (
	"mpc/internal/service"
	"mpc/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	BaseHandler
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	if assetService == nil {
		panic("assetService is required")
	}
	return &AssetHandler{
		BaseHandler:  NewBaseHandler(),
		assetService: assetService,
	}
}

// GetChains godoc
// @Summary      Get chains
// @Description  Get all chains
// @Tags         assets
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Response{payload=[]model.ChainResponse}
// @Failure      404  {object}  model.ErrorResponse
// @Router       /assets/chains [get]
func (h *AssetHandler) GetChains(c *gin.Context) {
	chains, err := h.assetService.GetChains(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, chains)
}

// GetTokensByChainID godoc
// @Summary      Get tokens by chain ID
// @Description  Get all tokens by chain ID
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        chain_id path int true "Chain ID"
// @Success      200  {object}  model.Response{payload=[]model.TokenResponse}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /assets/chains/{chain_id}/tokens [get]
func (h *AssetHandler) GetTokensByChainID(c *gin.Context) {
	chainID, err := h.parseChainID(c)
	if err != nil {
		c.Error(err)
		return
	}

	tokens, err := h.assetService.GetTokensByChainID(c.Request.Context(), chainID)
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, tokens)
}

// GetTokenBySymbol godoc
// @Summary      Get token by symbol
// @Description  Get token by symbol
// @Tags         assets
// @Accept       json
// @Produce      json
// @Param        chain_id path int true "Chain ID"
// @Param        symbol path string true "Symbol"
// @Success      200  {object}  model.Response{payload=model.Token}
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /assets/chains/{chain_id}/tokens/{symbol} [get]
func (h *AssetHandler) GetTokenBySymbol(c *gin.Context) {
	chainID, err := h.parseChainID(c)
	if err != nil {
		c.Error(err)
		return
	}

	symbol := c.Param("symbol")
	if symbol == "" {
		c.Error(errors.ErrInvalidSymbol)
		return
	}

	token, err := h.assetService.GetTokenBySymbol(c.Request.Context(), chainID, symbol)
	if err != nil {
		c.Error(err)
		return
	}
	h.SuccessResponse(c, token)
}

// Helper methods
func (h *AssetHandler) parseChainID(c *gin.Context) (int, error) {
	chainID := c.Param("chain_id")
	if chainID == "" {
		return 0, errors.ErrInvalidChainID
	}

	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		return 0, errors.ErrInvalidChainID
	}

	return chainIDInt, nil
}
