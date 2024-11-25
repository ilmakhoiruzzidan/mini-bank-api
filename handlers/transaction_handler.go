package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/models"
	"mini-bank-api/services"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
			"error":   err.Error()})
		return
	}

	sender, exists := c.Get("customer")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	customer, ok := sender.(*models.Customer)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to retrieve customer",
		})
		return
	}

	senderID := customer.ID
	transactionId, err := handler.transactionService.ProcessTransaction(senderID, paymentRequest.MerchantID, paymentRequest.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Payment successful",
		"data": gin.H{
			"transaction_id": transactionId,
			"sender_id":      senderID,
			"merchant_id":    paymentRequest.MerchantID,
			"amount":         paymentRequest.Amount,
			"timestamp":      time.Now(),
		},
	})
}
