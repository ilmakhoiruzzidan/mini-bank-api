package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`  // Omit if no data
	Error   string      `json:"error,omitempty"` // Omit if no error
}

func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

func LoginResponse(c *gin.Context, accessToken interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":      http.StatusOK,
		"accessToken": accessToken,
		"message":     message,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errorDetail string) {
	response := Response{
		Status:  statusCode,
		Message: message,
		Error:   errorDetail,
	}
	c.JSON(statusCode, response)
}
