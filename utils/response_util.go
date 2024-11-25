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

type LoginResponseDTO struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	AccessToken interface{} `json:"access_token"`
}

type CustomerProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func CurrentUserResponse(c *gin.Context, data interface{}, message string) {
	profile, ok := data.(CustomerProfile)
	if !ok {
		ErrorResponse(c, http.StatusInternalServerError, "Invalid data format", "Expected CustomerProfile format")
		return
	}
	response := Response{
		Status:  http.StatusOK,
		Message: message,
		Data:    profile,
	}
	c.JSON(http.StatusOK, response)
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
	c.JSON(http.StatusOK, LoginResponseDTO{
		Status:      http.StatusOK,
		Message:     message,
		AccessToken: accessToken,
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
