package handler

import (
	"mpc/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	BaseHandler
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
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
		h.HandleError(c, err)
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
// @Failure      404  {object}  model.ErrorResponse
// @Router       /assets/chains/{chain_id}/tokens [get]
func (h *AssetHandler) GetTokensByChainID(c *gin.Context) {
	chainID := c.Param("chain_id")
	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	tokens, err := h.assetService.GetTokensByChainID(c.Request.Context(), chainIDInt)
	if err != nil {
		h.HandleError(c, err)
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
// @Failure      404  {object}  model.ErrorResponse
// @Router       /assets/chains/{chain_id}/tokens/{symbol} [get]
func (h *AssetHandler) GetTokenBySymbol(c *gin.Context) {
	chainID := c.Param("chain_id")
	symbol := c.Param("symbol")
	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	token, err := h.assetService.GetTokenBySymbol(c.Request.Context(), chainIDInt, symbol)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	h.SuccessResponse(c, token)
}
