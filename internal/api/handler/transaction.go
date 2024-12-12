package handler

import (
	"mpc/internal/model"
	"mpc/internal/service"
	"mpc/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	BaseHandler
	txnService *service.TransactionService
}

func NewTransactionHandler(txnService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		BaseHandler: NewBaseHandler(),
		txnService:  txnService,
	}
}

// GetTransactions godoc
// @Summary      Get transactions
// @Description  Get all transactions by wallet address and chain ID(optional)
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        chain_id query string false "Chain ID"
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Success      200  {object}  model.Response{payload=model.TransactionListResponse}
// @Failure      400  {object}  model.ErrorResponse
// @Router       /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := h.GetUserID(c)

	chainID, _ := strconv.Atoi(c.DefaultQuery("chain_id", "11155111"))
	walletAddress := c.Query("wallet_address")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	res, err := h.txnService.GetTransactions(c.Request.Context(), userID, chainID, walletAddress, page, pageSize)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	h.SuccessResponse(c, res)
}

// CreateAndSubmitTransaction godoc
// @Summary      Create and submit transaction
// @Description  Create and submit transaction
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body model.CreateAndSubmitTransactionRequest true "Transaction request"
// @Success      200  {object}  model.Response{payload=model.Transaction}
// @Failure      400  {object}  model.ErrorResponse
// @Router       /transactions [post]
func (h *TransactionHandler) CreateAndSubmitTransaction(c *gin.Context) {
	userID := h.GetUserID(c)

	var req model.CreateAndSubmitTransactionRequest
	if err := utils.ValidateBody(c, &req); err != nil {
		h.HandleError(c, err)
		return
	}

	res, err := h.txnService.CreateAndSubmitTransaction(c.Request.Context(), userID, req)
	if err != nil {
		h.HandleError(c, err)
		return
	}
	h.SuccessResponse(c, res)
}
