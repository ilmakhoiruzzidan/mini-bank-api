package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/services"
	"net/http"
)

type CustomerHandler struct {
	Service services.CustomerServiceInterface
}

func NewCustomerHandler(service services.CustomerServiceInterface) *CustomerHandler {
	return &CustomerHandler{
		Service: service,
	}
}

func (ch *CustomerHandler) GetAllCustomers(c *gin.Context) {
	customers, err := ch.Service.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": customers})
}
