package handlers

import (
	"github.com/gin-gonic/gin"
	"mini-bank-api/services"
	"net/http"
)

type AuthHandler struct {
	Service services.AuthServiceInterface
}

func NewAuthHandler(service services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	token, err := handler.Service.Login(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      http.StatusOK,
		"message":     "login success",
		"accessToken": token,
	})
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	//handler.Service.Logout
}
