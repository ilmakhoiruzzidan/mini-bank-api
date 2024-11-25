package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/models"
	"mini-bank-api/services"
	"mini-bank-api/utils"
	"net/http"
	"time"
)

type TransactionHandler struct {
	transactionService services.TransactionServiceInterface
}

func NewTransactionHandler(transactionService services.TransactionServiceInterface) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (handler *TransactionHandler) CreateTransaction(c *gin.Context) {
	var paymentRequest struct {
		MerchantID string  `json:"merchant_id"`
		Amount     float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	sender, exists := c.Get("customer")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "")
		return
	}

	customer, ok := sender.(*models.Customer)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve customer", "")
		return
	}

	senderID := customer.ID
	transactionId, err := handler.transactionService.ProcessTransaction(senderID, paymentRequest.MerchantID, paymentRequest.Amount)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	utils.SuccessResponse(c,
		gin.H{
			"transaction_id": transactionId,
			"sender_id":      senderID,
			"merchant_id":    paymentRequest.MerchantID,
			"amount":         paymentRequest.Amount,
			"created_at":     time.Now(),
		}, "Payment successful")
}
